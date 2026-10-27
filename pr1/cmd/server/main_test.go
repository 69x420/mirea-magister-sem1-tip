package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHello(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got, want := rec.Body.String(), "Hello, world!\n"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestUserReturnsValidJSONWithUUID(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/user", nil)

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q", ct)
	}

	var got user
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, err := uuid.Parse(got.ID); err != nil {
		t.Fatalf("id %q is not a UUID: %v", got.ID, err)
	}
	if got.Name != "Gopher" {
		t.Fatalf("name = %q", got.Name)
	}
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	newRouter().ServeHTTP(rec, req)

	var got health
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Status != "ok" {
		t.Fatalf("status field = %q", got.Status)
	}
	if _, err := time.Parse(time.RFC3339, got.Time); err != nil {
		t.Fatalf("time %q is not RFC3339: %v", got.Time, err)
	}
}

func TestUserRejectsPost(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/user", nil)

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
