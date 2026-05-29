package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/thedevrems/tuto_lua/server/internal/httpx"
	"github.com/thedevrems/tuto_lua/server/internal/models"
	"github.com/thedevrems/tuto_lua/server/internal/store"
)

// CourseStore is the read surface the public course endpoints need.
type CourseStore interface {
	ListCourses(publishedOnly bool) ([]models.Course, error)
	GetCourseTree(slug string) (models.Course, error)
}

// CourseHandler serves the public course catalogue and content tree.
type CourseHandler struct {
	store CourseStore
}

// NewCourseHandler builds the handler around the course store.
func NewCourseHandler(s CourseStore) *CourseHandler {
	return &CourseHandler{store: s}
}

// List returns the published course catalogue (headers only, no content).
func (h *CourseHandler) List(w http.ResponseWriter, r *http.Request) {
	courses, err := h.store.ListCourses(true)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger les cours")
		return
	}
	httpx.JSON(w, http.StatusOK, courses)
}

// Tree returns one course with all its chapters, lessons, exercises and tests.
func (h *CourseHandler) Tree(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	course, err := h.store.GetCourseTree(slug)
	if errors.Is(err, store.ErrNotFound) {
		httpx.Error(w, http.StatusNotFound, "cours introuvable")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger le cours")
		return
	}
	httpx.JSON(w, http.StatusOK, course)
}
