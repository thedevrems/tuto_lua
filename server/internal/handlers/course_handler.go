package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/thedevrems/tuto_lua/server/internal/auth"
	"github.com/thedevrems/tuto_lua/server/internal/httpx"
	"github.com/thedevrems/tuto_lua/server/internal/models"
	"github.com/thedevrems/tuto_lua/server/internal/store"
)

// CourseStore is the read surface the public course endpoints need.
type CourseStore interface {
	ListCourses(publishedOnly bool) ([]models.Course, error)
	GetCourseBySlug(slug string) (models.Course, error)
	GetCourseTree(slug string) (models.Course, error)
	IsEnrolled(userID, courseID string) (bool, error)
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

// Tree returns one course with all its content, gated by access: free courses
// are open to everyone; paid ones require enrollment (or an admin).
func (h *CourseHandler) Tree(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	header, err := h.store.GetCourseBySlug(slug)
	if errors.Is(err, store.ErrNotFound) {
		httpx.Error(w, http.StatusNotFound, "cours introuvable")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger le cours")
		return
	}
	if !h.canAccess(r, header) {
		httpx.Error(w, http.StatusForbidden, "achetez ce cours pour accéder à son contenu")
		return
	}
	course, err := h.store.GetCourseTree(slug)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger le cours")
		return
	}
	httpx.JSON(w, http.StatusOK, course)
}

// canAccess applies the access policy for a course's full content.
func (h *CourseHandler) canAccess(r *http.Request, course models.Course) bool {
	if course.PriceCents == 0 {
		return true // free course
	}
	user, ok := auth.UserFrom(r.Context())
	if !ok {
		return false
	}
	if user.IsAdmin() {
		return true
	}
	enrolled, _ := h.store.IsEnrolled(user.ID, course.ID)
	return enrolled
}
