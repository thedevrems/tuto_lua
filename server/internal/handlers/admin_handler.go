package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/thedevrems/tuto_lua/server/internal/auth"
	"github.com/thedevrems/tuto_lua/server/internal/httpx"
	"github.com/thedevrems/tuto_lua/server/internal/models"
	"github.com/thedevrems/tuto_lua/server/internal/store"
)

// AdminHandler exposes the admin-only management endpoints. It uses the concrete
// store because administration touches many repositories at once.
type AdminHandler struct {
	store *store.Store
}

// NewAdminHandler builds the admin handler around the data store.
func NewAdminHandler(s *store.Store) *AdminHandler {
	return &AdminHandler{store: s}
}

// ListUsers returns every account (password hashes are never serialized).
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.ListUsers()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger les utilisateurs")
		return
	}
	if users == nil {
		users = []models.User{}
	}
	httpx.JSON(w, http.StatusOK, users)
}

// ListCourses returns every course, including unpublished drafts (admin view).
func (h *AdminHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.store.ListCourses(false)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger les cours")
		return
	}
	if courses == nil {
		courses = []models.Course{}
	}
	httpx.JSON(w, http.StatusOK, courses)
}

type grantAccessRequest struct {
	UserID   string `json:"userId"`
	CourseID string `json:"courseId"`
}

// GrantAccess enrolls a user in a course on behalf of the current admin.
func (h *AdminHandler) GrantAccess(w http.ResponseWriter, r *http.Request) {
	admin, _ := auth.UserFrom(r.Context())
	var req grantAccessRequest
	if !bind(w, r, &req) {
		return
	}
	if req.UserID == "" || req.CourseID == "" {
		httpx.Error(w, http.StatusBadRequest, "userId et courseId sont requis")
		return
	}
	adminID := admin.ID
	if err := h.store.CreateEnrollment(req.UserID, req.CourseID, "admin", &adminID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "octroi d'accès impossible (utilisateur ou cours introuvable ?)")
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}

// UserProgress returns a user's latest code ("last push") for every exercise.
func (h *AdminHandler) UserProgress(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	items, err := h.store.ListProgress(userID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger la progression")
		return
	}
	if items == nil {
		items = []models.Progress{}
	}
	httpx.JSON(w, http.StatusOK, items)
}

// UserCourses returns the courses a given user can access (free + unlocked).
func (h *AdminHandler) UserCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.store.ListUserCourses(chi.URLParam(r, "userId"))
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger les cours")
		return
	}
	if courses == nil {
		courses = []models.Course{}
	}
	httpx.JSON(w, http.StatusOK, courses)
}
