// Package handlers contains the HTTP endpoints; each method maps one route to
// one piece of behaviour and delegates the real work to services/stores.
package handlers

import (
	"errors"
	"net/http"

	"github.com/thedevrems/tuto_lua/server/internal/auth"
	"github.com/thedevrems/tuto_lua/server/internal/httpx"
	"github.com/thedevrems/tuto_lua/server/internal/store"
	"github.com/thedevrems/tuto_lua/server/internal/validate"
)

// AuthHandler exposes registration, login and the current-user endpoint.
type AuthHandler struct {
	svc *auth.Service
}

// NewAuthHandler builds the handler around the auth service.
func NewAuthHandler(svc *auth.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

// Register creates a new account and returns the user with a session token.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := httpx.Decode(w, r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.svc.Register(req.Username, req.Email, req.Password)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, res)
}

// Login authenticates by email or username and returns a fresh token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httpx.Decode(w, r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.svc.Login(req.Identifier, req.Password)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, res)
}

// Me returns the authenticated user injected by the auth middleware.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFrom(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "authentification requise")
		return
	}
	httpx.JSON(w, http.StatusOK, user)
}

// writeAuthError maps service errors to the right HTTP status.
func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrConflict):
		httpx.Error(w, http.StatusConflict, "ce nom d'utilisateur ou cet e-mail est déjà utilisé")
	case isValidationError(err):
		httpx.Error(w, http.StatusBadRequest, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "erreur interne")
	}
}

// isValidationError reports whether err is one of the input rule violations.
func isValidationError(err error) bool {
	return errors.Is(err, validate.ErrUsername) ||
		errors.Is(err, validate.ErrEmail) ||
		errors.Is(err, validate.ErrPassword)
}
