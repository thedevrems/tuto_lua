package handlers

import (
	"errors"
	"net/http"

	"github.com/thedevrems/tuto_lua/server/internal/auth"
	"github.com/thedevrems/tuto_lua/server/internal/httpx"
	"github.com/thedevrems/tuto_lua/server/internal/models"
)

// ProfileStore is the data surface the profile endpoints need.
type ProfileStore interface {
	ListUserCourses(userID string) ([]models.Course, error)
}

// ProfileHandler serves the authenticated user's own profile data.
type ProfileHandler struct {
	auth  *auth.Service
	store ProfileStore
}

// NewProfileHandler wires the auth service (password) and store (courses).
func NewProfileHandler(a *auth.Service, s ProfileStore) *ProfileHandler {
	return &ProfileHandler{auth: a, store: s}
}

// MyCourses returns the courses the current user can access (free + unlocked).
func (h *ProfileHandler) MyCourses(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFrom(r.Context())
	courses, err := h.store.ListUserCourses(user.ID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger vos cours")
		return
	}
	if courses == nil {
		courses = []models.Course{}
	}
	httpx.JSON(w, http.StatusOK, courses)
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// ChangePassword updates the current user's password after verifying the old one.
func (h *ProfileHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFrom(r.Context())
	var req changePasswordRequest
	if !bind(w, r, &req) {
		return
	}
	switch err := h.auth.ChangePassword(user.ID, req.CurrentPassword, req.NewPassword); {
	case errors.Is(err, auth.ErrInvalidCredentials):
		httpx.Error(w, http.StatusUnauthorized, "mot de passe actuel incorrect")
	case isValidationError(err):
		httpx.Error(w, http.StatusBadRequest, err.Error())
	case err != nil:
		httpx.Error(w, http.StatusInternalServerError, "impossible de changer le mot de passe")
	default:
		httpx.JSON(w, http.StatusNoContent, nil)
	}
}
