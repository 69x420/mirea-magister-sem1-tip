package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/zaldk/mirea-sem1-tid-private/pr3/internal/storage"
)

const (
	minTitleLen = 3
	maxTitleLen = 140
)

type Handlers struct {
	Store *storage.MemoryStore
}

func NewHandlers(store *storage.MemoryStore) *Handlers {
	return &Handlers{Store: store}
}

func (h *Handlers) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", h.GetTask)
	mux.HandleFunc("PATCH /tasks/{id}", h.PatchTask)
	mux.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)
	return mux
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	JSON(w, http.StatusOK, h.Store.List(q))
}

type createTaskRequest struct {
	Title string `json:"title"`
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && !strings.Contains(ct, "application/json") {
		BadRequest(w, "Content-Type must be application/json")
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid json: "+err.Error())
		return
	}

	title, ok := validateTitle(w, req.Title)
	if !ok {
		return
	}

	JSON(w, http.StatusCreated, h.Store.Create(title))
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	t, err := h.Store.Get(id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	JSON(w, http.StatusOK, t)
}

type patchTaskRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func (h *Handlers) PatchTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req patchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid json: "+err.Error())
		return
	}
	if req.Title == nil && req.Done == nil {
		BadRequest(w, "nothing to update: send title or done")
		return
	}

	if req.Title != nil {
		title, valid := validateTitle(w, *req.Title)
		if !valid {
			return
		}
		req.Title = &title
	}

	t, err := h.Store.Update(id, req.Title, req.Done)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	JSON(w, http.StatusOK, t)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.Store.Delete(id); err != nil {
		writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		BadRequest(w, "invalid id")
		return 0, false
	}
	return id, true
}

func validateTitle(w http.ResponseWriter, raw string) (string, bool) {
	title := strings.TrimSpace(raw)
	if title == "" {
		BadRequest(w, "title is required")
		return "", false
	}
	if n := len([]rune(title)); n < minTitleLen || n > maxTitleLen {
		Unprocessable(w, "title length must be between 3 and 140 characters")
		return "", false
	}
	return title, true
}

func writeStoreError(w http.ResponseWriter, err error) {
	if errors.Is(err, storage.ErrNotFound) {
		NotFound(w, "task not found")
		return
	}
	Internal(w, "unexpected error")
}
