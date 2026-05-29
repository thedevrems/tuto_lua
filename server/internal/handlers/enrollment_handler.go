package handlers

import (
	"net/http"

	"github.com/thedevrems/tuto_lua/server/internal/auth"
	"github.com/thedevrems/tuto_lua/server/internal/httpx"
)

// EnrollmentStore is the data surface the enrollment endpoint needs.
type EnrollmentStore interface {
	ListEnrollmentCourseIDs(userID string) ([]string, error)
}

// EnrollmentHandler reports which courses the current user can access.
type EnrollmentHandler struct {
	store EnrollmentStore
}

// NewEnrollmentHandler builds the handler around the enrollment store.
func NewEnrollmentHandler(s EnrollmentStore) *EnrollmentHandler {
	return &EnrollmentHandler{store: s}
}

// Mine returns the ids of the courses the authenticated user is enrolled in.
func (h *EnrollmentHandler) Mine(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFrom(r.Context())
	ids, err := h.store.ListEnrollmentCourseIDs(user.ID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger les accès")
		return
	}
	if ids == nil {
		ids = []string{}
	}
	httpx.JSON(w, http.StatusOK, ids)
}
