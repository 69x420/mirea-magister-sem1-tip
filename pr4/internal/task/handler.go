package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

const (
	minTitleLen = 3
	maxTitleLen = 100
	defaultPage = 1
	defaultSize = 10
	maxPageSize = 100
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Get("/{id}", h.get)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tasks := h.repo.List()

	if q := strings.TrimSpace(r.URL.Query().Get("q")); q != "" {
		tasks = filterByTitle(tasks, q)
	}
	if raw := r.URL.Query().Get("done"); raw != "" {
		done, err := strconv.ParseBool(raw)
		if err != nil {
			httpError(w, http.StatusBadRequest, "done must be true or false")
			return
		}
		tasks = filterByDone(tasks, done)
	}

	page, err := queryInt(r, "page", defaultPage)
	if err != nil || page < 1 {
		httpError(w, http.StatusBadRequest, "page must be a positive integer")
		return
	}
	limit, err := queryInt(r, "limit", defaultSize)
	if err != nil || limit < 1 || limit > maxPageSize {
		httpError(w, http.StatusBadRequest, "limit must be between 1 and 100")
		return
	}

	total := len(tasks)
	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	writeJSON(w, http.StatusOK, listResponse{
		Items: tasks[start:end],
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	t, err := h.repo.Get(id)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

type createReq struct {
	Title string `json:"title"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && !strings.Contains(ct, "application/json") {
		httpError(w, http.StatusBadRequest, "Content-Type must be application/json")
		return
	}

	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	title, ok := validateTitle(w, req.Title)
	if !ok {
		return
	}
	writeJSON(w, http.StatusCreated, h.repo.Create(title))
}

type updateReq struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req updateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}
	if req.Title == nil && req.Done == nil {
		httpError(w, http.StatusBadRequest, "nothing to update: send title or done")
		return
	}
	if req.Title != nil {
		title, valid := validateTitle(w, *req.Title)
		if !valid {
			return
		}
		req.Title = &title
	}

	t, err := h.repo.Update(id, req.Title, req.Done)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.repo.Delete(id); err != nil {
		writeRepoError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type listResponse struct {
	Items []*Task `json:"items"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, errorResponse{Error: msg})
}

func writeRepoError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotFound) {
		httpError(w, http.StatusNotFound, "task not found")
		return
	}
	httpError(w, http.StatusInternalServerError, "unexpected error")
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		httpError(w, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}

func validateTitle(w http.ResponseWriter, raw string) (string, bool) {
	title := strings.TrimSpace(raw)
	if title == "" {
		httpError(w, http.StatusBadRequest, "title is required")
		return "", false
	}
	if n := len([]rune(title)); n < minTitleLen || n > maxTitleLen {
		httpError(w, http.StatusUnprocessableEntity, "title length must be between 3 and 100 characters")
		return "", false
	}
	return title, true
}

func queryInt(r *http.Request, name string, fallback int) (int, error) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return fallback, nil
	}
	return strconv.Atoi(raw)
}

func filterByTitle(tasks []*Task, q string) []*Task {
	q = strings.ToLower(q)
	out := make([]*Task, 0, len(tasks))
	for _, t := range tasks {
		if strings.Contains(strings.ToLower(t.Title), q) {
			out = append(out, t)
		}
	}
	return out
}

func filterByDone(tasks []*Task, done bool) []*Task {
	out := make([]*Task, 0, len(tasks))
	for _, t := range tasks {
		if t.Done == done {
			out = append(out, t)
		}
	}
	return out
}
