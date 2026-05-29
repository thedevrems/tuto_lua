package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/thedevrems/tuto_lua/server/internal/auth"
	"github.com/thedevrems/tuto_lua/server/internal/httpx"
	"github.com/thedevrems/tuto_lua/server/internal/models"
)

// ProgressStore is the data surface the progress endpoints need.
type ProgressStore interface {
	ListProgress(userID string) ([]models.Progress, error)
	UpsertProgress(userID, exerciseID, code string, completed bool) (models.Progress, error)
}

// ProgressHandler reads and persists a user's per-exercise progress.
type ProgressHandler struct {
	store ProgressStore
}

// NewProgressHandler builds the handler around the progress store.
func NewProgressHandler(s ProgressStore) *ProgressHandler {
	return &ProgressHandler{store: s}
}

type saveProgressRequest struct {
	Code      string `json:"code"`
	Completed bool   `json:"completed"`
}

// List returns all progress rows for the authenticated user.
func (h *ProgressHandler) List(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFrom(r.Context())
	items, err := h.store.ListProgress(user.ID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger la progression")
		return
	}
	if items == nil {
		items = []models.Progress{}
	}
	httpx.JSON(w, http.StatusOK, items)
}

// Save upserts the user's latest code and completion state for one exercise.
func (h *ProgressHandler) Save(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFrom(r.Context())
	exerciseID := chi.URLParam(r, "exerciseId")

	var req saveProgressRequest
	if err := httpx.Decode(w, r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.store.UpsertProgress(user.ID, exerciseID, req.Code, req.Completed)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible d'enregistrer la progression")
		return
	}
	httpx.JSON(w, http.StatusOK, saved)
}
