package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func do(t *testing.T, path string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	return rec
}

func TestRootPlainText(t *testing.T) {
	rec := do(t, "/", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Fatalf("Content-Type = %q", ct)
	}
	if got, want := rec.Body.String(), "Hello, Go project structure!\n"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestPingJSON(t *testing.T) {
	rec := do(t, "/ping", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	var got struct {
		Status string `json:"status"`
		Time   string `json:"time"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Status != "ok" {
		t.Fatalf("status field = %q", got.Status)
	}
	if _, err := time.Parse(time.RFC3339, got.Time); err != nil {
		t.Fatalf("time %q: %v", got.Time, err)
	}
}

func TestRequestIDGenerated(t *testing.T) {
	rec := do(t, "/ping", nil)

	if id := rec.Header().Get("X-Request-Id"); len(id) != 16 {
		t.Fatalf("generated X-Request-Id = %q, want 16 hex chars", id)
	}
}

func TestRequestIDPropagated(t *testing.T) {
	rec := do(t, "/ping", map[string]string{"X-Request-Id": "demo-123"})

	if id := rec.Header().Get("X-Request-Id"); id != "demo-123" {
		t.Fatalf("X-Request-Id = %q, want demo-123", id)
	}
}

func TestFailUsesUnifiedErrorFormat(t *testing.T) {
	rec := do(t, "/fail", nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}

	var got struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "bad_request_example" {
		t.Fatalf("error = %q", got.Error)
	}
}

func TestUnknownRouteIs404(t *testing.T) {
	rec := do(t, "/nope", nil)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestPortFromEnv(t *testing.T) {
	t.Setenv("PORT", "9091")
	if got := Port(); got != "9091" {
		t.Fatalf("Port() = %q, want 9091", got)
	}
}

func TestPortDefault(t *testing.T) {
	t.Setenv("PORT", "")
	if got := Port(); got != "8080" {
		t.Fatalf("Port() = %q, want 8080", got)
	}
}
