package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/zaldk/mirea-sem1-tid-private/pr6/internal/models"
)

type Handlers struct{ db *gorm.DB }

func NewHandlers(db *gorm.DB) *Handlers { return &Handlers{db: db} }

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type createUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var in createUserReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil ||
		strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Email) == "" {
		writeErr(w, http.StatusBadRequest, "name and email are required")
		return
	}

	u := models.User{Name: strings.TrimSpace(in.Name), Email: strings.TrimSpace(in.Email)}
	if err := h.db.Create(&u).Error; err != nil {
		writeErr(w, http.StatusConflict, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, u)
}

func (h *Handlers) ListUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := h.db.Preload("Notes").Order("id").Find(&users).Error; err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, users)
}

type createNoteReq struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	UserID  uint     `json:"userId"`
	Tags    []string `json:"tags"`
}

func (h *Handlers) CreateNote(w http.ResponseWriter, r *http.Request) {
	var in createNoteReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil ||
		strings.TrimSpace(in.Title) == "" || in.UserID == 0 {
		writeErr(w, http.StatusBadRequest, "title and userId are required")
		return
	}

	var author models.User
	if err := h.db.First(&author, in.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeErr(w, http.StatusNotFound, "user not found")
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	tags := make([]models.Tag, 0, len(in.Tags))
	for _, name := range in.Tags {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		tag := models.Tag{Name: name}
		if err := h.db.Where(models.Tag{Name: name}).FirstOrCreate(&tag).Error; err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		tags = append(tags, tag)
	}

	note := models.Note{
		Title:   strings.TrimSpace(in.Title),
		Content: in.Content,
		UserID:  in.UserID,
		Tags:    tags,
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&note).Error; err != nil {
			return err
		}
		return tx.Preload("User").Preload("Tags").First(&note, note.ID).Error
	})
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, note)
}

func (h *Handlers) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}

	var note models.Note
	if err := h.db.Preload("User").Preload("Tags").First(&note, id).Error; err != nil {
		writeErr(w, http.StatusNotFound, "note not found")
		return
	}

	writeJSON(w, http.StatusOK, note)
}

func (h *Handlers) ListNotes(w http.ResponseWriter, r *http.Request) {
	var notes []models.Note
	if err := h.db.Preload("User").Preload("Tags").Order("id").Find(&notes).Error; err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, notes)
}
