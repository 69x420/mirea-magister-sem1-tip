package task

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newServer() http.Handler {
	return NewHandler(NewRepo()).Routes()
}

func do(t *testing.T, srv http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

func decodeTask(t *testing.T, rec *httptest.ResponseRecorder) Task {
	t.Helper()
	var task Task
	if err := json.NewDecoder(rec.Body).Decode(&task); err != nil {
		t.Fatalf("decode task: %v", err)
	}
	return task
}

func decodeList(t *testing.T, rec *httptest.ResponseRecorder) listResponse {
	t.Helper()
	var list listResponse
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	return list
}

func seed(t *testing.T, srv http.Handler, titles ...string) {
	t.Helper()
	for _, title := range titles {
		if rec := do(t, srv, http.MethodPost, "/", `{"title":"`+title+`"}`); rec.Code != http.StatusCreated {
			t.Fatalf("seed %q: %d %s", title, rec.Code, rec.Body)
		}
	}
}

func TestCreate(t *testing.T) {
	srv := newServer()
	rec := do(t, srv, http.MethodPost, "/", `{"title":"  Выучить chi  "}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	task := decodeTask(t, rec)
	if task.ID != 1 || task.Title != "Выучить chi" || task.Done {
		t.Fatalf("unexpected task: %+v", task)
	}
}

func TestCreateValidation(t *testing.T) {
	srv := newServer()
	cases := []struct {
		name string
		body string
		code int
	}{
		{"empty", `{"title":"   "}`, http.StatusBadRequest},
		{"short", `{"title":"ab"}`, http.StatusUnprocessableEntity},
		{"long", `{"title":"` + strings.Repeat("я", 101) + `"}`, http.StatusUnprocessableEntity},
		{"broken json", `{"title":`, http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if rec := do(t, srv, http.MethodPost, "/", tc.body); rec.Code != tc.code {
				t.Fatalf("status = %d, want %d", rec.Code, tc.code)
			}
		})
	}
}

func TestGet(t *testing.T) {
	srv := newServer()
	seed(t, srv, "Выучить chi")

	rec := do(t, srv, http.MethodGet, "/1", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if decodeTask(t, rec).Title != "Выучить chi" {
		t.Fatalf("wrong title")
	}

	if rec := do(t, srv, http.MethodGet, "/9999", ""); rec.Code != http.StatusNotFound {
		t.Fatalf("missing task status = %d", rec.Code)
	}
	if rec := do(t, srv, http.MethodGet, "/abc", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad id status = %d", rec.Code)
	}
}

func TestUpdate(t *testing.T) {
	srv := newServer()
	seed(t, srv, "Выучить chi")

	rec := do(t, srv, http.MethodPut, "/1", `{"title":"Выучить chi глубже","done":true}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	task := decodeTask(t, rec)
	if task.Title != "Выучить chi глубже" || !task.Done {
		t.Fatalf("unexpected task: %+v", task)
	}

	if rec := do(t, srv, http.MethodPut, "/1", `{}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("empty update status = %d", rec.Code)
	}
	if rec := do(t, srv, http.MethodPut, "/9999", `{"done":true}`); rec.Code != http.StatusNotFound {
		t.Fatalf("missing update status = %d", rec.Code)
	}
}

func TestDelete(t *testing.T) {
	srv := newServer()
	seed(t, srv, "Выучить chi")

	rec := do(t, srv, http.MethodDelete, "/1", "")
	if rec.Code != http.StatusNoContent || rec.Body.Len() != 0 {
		t.Fatalf("status = %d body = %q", rec.Code, rec.Body)
	}
	if again := do(t, srv, http.MethodDelete, "/1", ""); again.Code != http.StatusNotFound {
		t.Fatalf("second delete = %d", again.Code)
	}
}

func TestFilterByDone(t *testing.T) {
	srv := newServer()
	seed(t, srv, "Первая задача", "Вторая задача", "Третья задача")
	do(t, srv, http.MethodPut, "/2", `{"done":true}`)

	rec := do(t, srv, http.MethodGet, "/?done=true", "")
	list := decodeList(t, rec)
	if list.Total != 1 || len(list.Items) != 1 || list.Items[0].ID != 2 {
		t.Fatalf("done filter: %+v", list)
	}

	if rec := do(t, srv, http.MethodGet, "/?done=maybe", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad done status = %d", rec.Code)
	}
}

func TestPagination(t *testing.T) {
	srv := newServer()
	seed(t, srv, "Задача один", "Задача два", "Задача три", "Задача четыре", "Задача пять")

	list := decodeList(t, do(t, srv, http.MethodGet, "/?page=2&limit=2", ""))
	if list.Total != 5 || list.Page != 2 || list.Limit != 2 || len(list.Items) != 2 {
		t.Fatalf("page 2: %+v", list)
	}
	if list.Items[0].ID != 3 {
		t.Fatalf("first item on page 2 = %d", list.Items[0].ID)
	}

	last := decodeList(t, do(t, srv, http.MethodGet, "/?page=3&limit=2", ""))
	if len(last.Items) != 1 || last.Items[0].ID != 5 {
		t.Fatalf("last page: %+v", last)
	}

	if rec := do(t, srv, http.MethodGet, "/?limit=0", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("limit=0 status = %d", rec.Code)
	}
	if rec := do(t, srv, http.MethodGet, "/?page=abc", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("page=abc status = %d", rec.Code)
	}
}

func TestSearchQuery(t *testing.T) {
	srv := newServer()
	seed(t, srv, "Купить молоко", "Купить хлеб", "Позвонить маме")

	list := decodeList(t, do(t, srv, http.MethodGet, "/?q=купить", ""))
	if list.Total != 2 {
		t.Fatalf("total = %d, want 2", list.Total)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	if rec := do(t, newServer(), http.MethodPatch, "/", `{"title":"Выучить chi"}`); rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}
