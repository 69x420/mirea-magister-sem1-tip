package httpapi

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/zaldk/mirea-sem1-tid-private/pr6/pkg/middleware"
)

func BuildRouter(d *gorm.DB) *chi.Mux {
	r := chi.NewRouter()
	h := NewHandlers(d)

	r.Use(middleware.Logger)
	r.Use(middleware.SimpleCORS)

	r.Get("/health", h.Health)
	r.Post("/users", h.CreateUser)
	r.Get("/users", h.ListUsers)
	r.Post("/notes", h.CreateNote)
	r.Get("/notes", h.ListNotes)
	r.Get("/notes/{id}", h.GetNoteByID)

	return r
}
