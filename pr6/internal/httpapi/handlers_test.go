package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zaldk/mirea-sem1-tid-private/pr6/internal/models"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_foreign_keys=1",
		filepath.Join(t.TempDir(), "test.db"))

	d, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	if err := d.AutoMigrate(&models.User{}, &models.Note{}, &models.Tag{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}
	return d
}

func do(t *testing.T, r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestHealth(t *testing.T) {
	r := BuildRouter(newTestDB(t))

	rec := do(t, r, http.MethodGet, "/health", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d", rec.Code)
	}
	if got := rec.Body.String(); got != "{\"status\":\"ok\"}\n" {
		t.Errorf("body = %q", got)
	}
}

func TestCreateUser(t *testing.T) {
	d := newTestDB(t)
	r := BuildRouter(d)

	rec := do(t, r, http.MethodPost, "/users", `{"name":"Alice","email":"alice@example.com"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("code = %d, body = %s", rec.Code, rec.Body)
	}

	var user models.User
	if err := json.Unmarshal(rec.Body.Bytes(), &user); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if user.ID != 1 || user.Name != "Alice" || user.Email != "alice@example.com" {
		t.Errorf("user = %+v", user)
	}

	var count int64
	d.Model(&models.User{}).Count(&count)
	if count != 1 {
		t.Errorf("count = %d", count)
	}
}

func TestCreateUserValidation(t *testing.T) {
	r := BuildRouter(newTestDB(t))

	for _, body := range []string{
		`{"name":"","email":"a@example.com"}`,
		`{"name":"Bob","email":""}`,
		`{}`,
		`not json`,
	} {
		rec := do(t, r, http.MethodPost, "/users", body)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("body %q: code = %d, want 400", body, rec.Code)
		}
	}
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	r := BuildRouter(newTestDB(t))

	body := `{"name":"Alice","email":"alice@example.com"}`
	if rec := do(t, r, http.MethodPost, "/users", body); rec.Code != http.StatusCreated {
		t.Fatalf("first create: %d", rec.Code)
	}

	rec := do(t, r, http.MethodPost, "/users", body)
	if rec.Code != http.StatusConflict {
		t.Fatalf("code = %d, want 409", rec.Code)
	}
}

func createUser(t *testing.T, r http.Handler) {
	t.Helper()

	rec := do(t, r, http.MethodPost, "/users", `{"name":"Alice","email":"alice@example.com"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create user: %d", rec.Code)
	}
}

func TestCreateNote(t *testing.T) {
	r := BuildRouter(newTestDB(t))
	createUser(t, r)

	body := `{"title":"Первая заметка","content":"Текст заметки","userId":1,"tags":["go","gorm"]}`
	rec := do(t, r, http.MethodPost, "/notes", body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("code = %d, body = %s", rec.Code, rec.Body)
	}

	var note models.Note
	if err := json.Unmarshal(rec.Body.Bytes(), &note); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if note.ID != 1 || note.Title != "Первая заметка" || note.UserID != 1 {
		t.Errorf("note = %+v", note)
	}
	if note.User.Name != "Alice" {
		t.Errorf("user = %+v", note.User)
	}
	if len(note.Tags) != 2 || note.Tags[0].Name != "go" || note.Tags[1].Name != "gorm" {
		t.Errorf("tags = %+v", note.Tags)
	}
}

func TestCreateNoteReusesTags(t *testing.T) {
	d := newTestDB(t)
	r := BuildRouter(d)
	createUser(t, r)

	first := `{"title":"Заметка 1","content":"a","userId":1,"tags":["go","gorm"]}`
	second := `{"title":"Заметка 2","content":"b","userId":1,"tags":["go","sqlite"]}`

	for _, body := range []string{first, second} {
		if rec := do(t, r, http.MethodPost, "/notes", body); rec.Code != http.StatusCreated {
			t.Fatalf("create note: %d, body = %s", rec.Code, rec.Body)
		}
	}

	var tags []models.Tag
	d.Order("id").Find(&tags)
	if len(tags) != 3 {
		t.Fatalf("tags = %+v, want 3 unique tags", tags)
	}
}

func TestCreateNoteValidation(t *testing.T) {
	r := BuildRouter(newTestDB(t))
	createUser(t, r)

	for _, body := range []string{
		`{"title":"","content":"x","userId":1}`,
		`{"title":"Заметка","content":"x","userId":0}`,
		`not json`,
	} {
		rec := do(t, r, http.MethodPost, "/notes", body)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("body %q: code = %d, want 400", body, rec.Code)
		}
	}
}

func TestCreateNoteUnknownUser(t *testing.T) {
	r := BuildRouter(newTestDB(t))

	rec := do(t, r, http.MethodPost, "/notes",
		`{"title":"Заметка","content":"x","userId":42}`)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("code = %d, want 404", rec.Code)
	}
}

func TestGetNoteByID(t *testing.T) {
	r := BuildRouter(newTestDB(t))
	createUser(t, r)

	body := `{"title":"Первая заметка","content":"Текст","userId":1,"tags":["go","gorm"]}`
	if rec := do(t, r, http.MethodPost, "/notes", body); rec.Code != http.StatusCreated {
		t.Fatalf("create note: %d", rec.Code)
	}

	rec := do(t, r, http.MethodGet, "/notes/1", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d", rec.Code)
	}

	var note models.Note
	if err := json.Unmarshal(rec.Body.Bytes(), &note); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if note.User.Name != "Alice" || len(note.Tags) != 2 {
		t.Errorf("note = %+v", note)
	}
}

func TestGetNoteErrors(t *testing.T) {
	r := BuildRouter(newTestDB(t))

	if rec := do(t, r, http.MethodGet, "/notes/abc", ""); rec.Code != http.StatusBadRequest {
		t.Errorf("abc: code = %d, want 400", rec.Code)
	}
	if rec := do(t, r, http.MethodGet, "/notes/0", ""); rec.Code != http.StatusBadRequest {
		t.Errorf("0: code = %d, want 400", rec.Code)
	}
	if rec := do(t, r, http.MethodGet, "/notes/99", ""); rec.Code != http.StatusNotFound {
		t.Errorf("99: code = %d, want 404", rec.Code)
	}
}

func TestListNotesWithPreload(t *testing.T) {
	r := BuildRouter(newTestDB(t))
	createUser(t, r)

	for _, title := range []string{"Заметка 1", "Заметка 2"} {
		body := fmt.Sprintf(`{"title":%q,"content":"x","userId":1,"tags":["go"]}`, title)
		if rec := do(t, r, http.MethodPost, "/notes", body); rec.Code != http.StatusCreated {
			t.Fatalf("create: %d", rec.Code)
		}
	}

	rec := do(t, r, http.MethodGet, "/notes", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d", rec.Code)
	}

	var notes []models.Note
	if err := json.Unmarshal(rec.Body.Bytes(), &notes); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("len = %d", len(notes))
	}
	for _, note := range notes {
		if note.User.Name != "Alice" || len(note.Tags) != 1 {
			t.Errorf("note %d not preloaded: %+v", note.ID, note)
		}
	}
}
