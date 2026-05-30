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

// NotificationStore is the data surface the notification endpoints need.
type NotificationStore interface {
	ListNotifications(userID string) ([]models.Notification, error)
	CountUnread(userID string) (int, error)
	MarkNotificationRead(userID, id string) error
	MarkAllNotificationsRead(userID string) error
}

// NotificationHandler serves the current user's in-app notifications.
type NotificationHandler struct {
	store NotificationStore
}

// NewNotificationHandler builds the handler around the notification store.
func NewNotificationHandler(s NotificationStore) *NotificationHandler {
	return &NotificationHandler{store: s}
}

// List returns the user's notifications plus the unread count.
func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFrom(r.Context())
	items, err := h.store.ListNotifications(user.ID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger les notifications")
		return
	}
	unread, _ := h.store.CountUnread(user.ID)
	if items == nil {
		items = []models.Notification{}
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"notifications": items, "unread": unread})
}

// MarkRead marks a single notification as read.
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFrom(r.Context())
	err := h.store.MarkNotificationRead(user.ID, chi.URLParam(r, "id"))
	if errors.Is(err, store.ErrNotFound) {
		httpx.Error(w, http.StatusNotFound, "notification introuvable")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "opération impossible")
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}

// MarkAllRead marks every notification of the user as read.
func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFrom(r.Context())
	if err := h.store.MarkAllNotificationsRead(user.ID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "opération impossible")
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}
