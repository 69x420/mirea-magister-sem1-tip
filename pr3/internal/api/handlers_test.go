package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zaldk/mirea-sem1-tid-private/pr3/internal/storage"
)

func newTestServer() http.Handler {
	return Logging(NewHandlers(storage.NewMemoryStore()).Routes())
}

func do(t *testing.T, srv http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

func decodeTask(t *testing.T, rec *httptest.ResponseRecorder) storage.Task {
	t.Helper()
	var task storage.Task
	if err := json.NewDecoder(rec.Body).Decode(&task); err != nil {
		t.Fatalf("decode task: %v", err)
	}
	return task
}

func TestHealth(t *testing.T) {
	rec := do(t, newTestServer(), http.MethodGet, "/health", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if got, want := strings.TrimSpace(rec.Body.String()), `{"status":"ok"}`; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestCreateAndGet(t *testing.T) {
	srv := newTestServer()

	created := do(t, srv, http.MethodPost, "/tasks", `{"title":"Купить молоко"}`)
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body = %s", created.Code, created.Body)
	}
	task := decodeTask(t, created)
	if task.ID != 1 || task.Title != "Купить молоко" || task.Done {
		t.Fatalf("unexpected task: %+v", task)
	}

	got := do(t, srv, http.MethodGet, "/tasks/1", "")
	if got.Code != http.StatusOK {
		t.Fatalf("get status = %d", got.Code)
	}
	if decodeTask(t, got).Title != "Купить молоко" {
		t.Fatalf("wrong title in GET")
	}
}

func TestCreateTrimsTitle(t *testing.T) {
	rec := do(t, newTestServer(), http.MethodPost, "/tasks", `{"title":"  Позвонить маме  "}`)

	if got := decodeTask(t, rec).Title; got != "Позвонить маме" {
		t.Fatalf("title = %q", got)
	}
}

func TestCreateValidation(t *testing.T) {
	srv := newTestServer()

	cases := []struct {
		name string
		body string
		code int
	}{
		{"empty title", `{"title":"   "}`, http.StatusBadRequest},
		{"too short", `{"title":"ab"}`, http.StatusUnprocessableEntity},
		{"too long", `{"title":"` + strings.Repeat("x", 141) + `"}`, http.StatusUnprocessableEntity},
		{"broken json", `{"title":`, http.StatusBadRequest},
		{"wrong content type", `{"title":"Нормальный заголовок"}`, http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(tc.body))
			if tc.name == "wrong content type" {
				req.Header.Set("Content-Type", "text/plain")
			} else {
				req.Header.Set("Content-Type", "application/json")
			}
			rec := httptest.NewRecorder()
			srv.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Fatalf("status = %d, want %d (body %s)", rec.Code, tc.code, rec.Body)
			}
		})
	}
}

func TestListAndFilter(t *testing.T) {
	srv := newTestServer()
	for _, title := range []string{"Купить молоко", "Купить хлеб", "Позвонить маме"} {
		if rec := do(t, srv, http.MethodPost, "/tasks", `{"title":"`+title+`"}`); rec.Code != http.StatusCreated {
			t.Fatalf("seed %q: %d", title, rec.Code)
		}
	}

	all := do(t, srv, http.MethodGet, "/tasks", "")
	var tasks []storage.Task
	if err := json.NewDecoder(all.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("len = %d, want 3", len(tasks))
	}

	filtered := do(t, srv, http.MethodGet, "/tasks?q=купить", "")
	tasks = nil
	if err := json.NewDecoder(filtered.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode filtered: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("filtered len = %d, want 2", len(tasks))
	}
}

func TestPatchDone(t *testing.T) {
	srv := newTestServer()
	do(t, srv, http.MethodPost, "/tasks", `{"title":"Купить молоко"}`)

	rec := do(t, srv, http.MethodPatch, "/tasks/1", `{"done":true}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	if task := decodeTask(t, rec); !task.Done {
		t.Fatalf("done = false, want true")
	}

	empty := do(t, srv, http.MethodPatch, "/tasks/1", `{}`)
	if empty.Code != http.StatusBadRequest {
		t.Fatalf("empty patch status = %d, want 400", empty.Code)
	}
}

func TestDelete(t *testing.T) {
	srv := newTestServer()
	do(t, srv, http.MethodPost, "/tasks", `{"title":"Купить молоко"}`)

	rec := do(t, srv, http.MethodDelete, "/tasks/1", "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("body = %q, want empty", rec.Body)
	}
	if again := do(t, srv, http.MethodDelete, "/tasks/1", ""); again.Code != http.StatusNotFound {
		t.Fatalf("second delete status = %d, want 404", again.Code)
	}
}

func TestErrors(t *testing.T) {
	srv := newTestServer()
	do(t, srv, http.MethodPost, "/tasks", `{"title":"Купить молоко"}`)

	cases := []struct {
		method string
		path   string
		code   int
	}{
		{http.MethodGet, "/tasks/abc", http.StatusBadRequest},
		{http.MethodGet, "/tasks/0", http.StatusBadRequest},
		{http.MethodGet, "/tasks/9999", http.StatusNotFound},
		{http.MethodGet, "/unknown", http.StatusNotFound},
		{http.MethodDelete, "/tasks/9999", http.StatusNotFound},
	}

	for _, tc := range cases {
		if rec := do(t, srv, tc.method, tc.path, ""); rec.Code != tc.code {
			t.Fatalf("%s %s = %d, want %d", tc.method, tc.path, rec.Code, tc.code)
		}
	}
}

func TestMethodNotAllowed(t *testing.T) {
	rec := do(t, newTestServer(), http.MethodPut, "/tasks", `{"title":"Купить молоко"}`)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}

func TestCORSHeaders(t *testing.T) {
	srv := CORS(newTestServer())

	rec := do(t, srv, http.MethodGet, "/health", "")
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("Allow-Origin = %q", got)
	}

	preflight := do(t, srv, http.MethodOptions, "/tasks", "")
	if preflight.Code != http.StatusNoContent {
		t.Fatalf("preflight status = %d", preflight.Code)
	}
}
