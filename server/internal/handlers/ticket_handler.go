package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/thedevrems/tuto_lua/server/internal/auth"
	"github.com/thedevrems/tuto_lua/server/internal/httpx"
	"github.com/thedevrems/tuto_lua/server/internal/store"
	"github.com/thedevrems/tuto_lua/server/internal/ticket"
)

// TicketHandler serves both the user and admin ticket (report/devis) endpoints.
type TicketHandler struct {
	svc *ticket.Service
}

// NewTicketHandler builds the handler around the ticket service.
func NewTicketHandler(svc *ticket.Service) *TicketHandler {
	return &TicketHandler{svc: svc}
}

type createReportRequest struct {
	Subject  string `json:"subject"`
	Category string `json:"category"`
	Body     string `json:"body"`
}

type messageRequest struct {
	Body string `json:"body"`
}

type addMemberRequest struct {
	UserID string `json:"userId"`
}

// Create opens a new support report for the current user.
func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	actor, _ := auth.UserFrom(r.Context())
	var req createReportRequest
	if !bind(w, r, &req) {
		return
	}
	t, err := h.svc.CreateReport(actor, req.Subject, req.Category, req.Body)
	if err != nil {
		writeTicketError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, t)
}

// ListMine returns the tickets the current user participates in.
func (h *TicketHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	actor, _ := auth.UserFrom(r.Context())
	list, err := h.svc.ListMine(actor.ID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger vos conversations")
		return
	}
	httpx.JSON(w, http.StatusOK, list)
}

// Get returns one ticket with its messages and members (if accessible).
func (h *TicketHandler) Get(w http.ResponseWriter, r *http.Request) {
	actor, _ := auth.UserFrom(r.Context())
	t, err := h.svc.Detail(actor, chi.URLParam(r, "id"))
	if err != nil {
		writeTicketError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, t)
}

// PostMessage adds a reply to a ticket conversation.
func (h *TicketHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
	actor, _ := auth.UserFrom(r.Context())
	var req messageRequest
	if !bind(w, r, &req) {
		return
	}
	msg, err := h.svc.PostMessage(actor, chi.URLParam(r, "id"), req.Body)
	if err != nil {
		writeTicketError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, msg)
}

// ListAll returns every ticket, optionally filtered by ?type= (admin only).
func (h *TicketHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListAll(r.URL.Query().Get("type"))
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "impossible de charger les conversations")
		return
	}
	httpx.JSON(w, http.StatusOK, list)
}

// Close closes a ticket and notifies its members (admin only).
func (h *TicketHandler) Close(w http.ResponseWriter, r *http.Request) {
	actor, _ := auth.UserFrom(r.Context())
	if err := h.svc.Close(actor, chi.URLParam(r, "id")); err != nil {
		writeTicketError(w, err)
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}

// AddMember adds another participant to a ticket (admin only).
func (h *TicketHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	actor, _ := auth.UserFrom(r.Context())
	var req addMemberRequest
	if !bind(w, r, &req) {
		return
	}
	if err := h.svc.AddMember(actor, chi.URLParam(r, "id"), req.UserID); err != nil {
		writeTicketError(w, err)
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}

// writeTicketError maps ticket service errors to HTTP statuses.
func writeTicketError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ticket.ErrForbidden):
		httpx.Error(w, http.StatusForbidden, err.Error())
	case errors.Is(err, ticket.ErrEmpty):
		httpx.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.Error(w, http.StatusNotFound, "conversation introuvable")
	default:
		httpx.Error(w, http.StatusInternalServerError, "opération impossible")
	}
}
