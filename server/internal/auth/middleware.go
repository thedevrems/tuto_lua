package auth

import (
	"net/http"
	"strings"

	"github.com/thedevrems/tuto_lua/server/internal/httpx"
	"github.com/thedevrems/tuto_lua/server/internal/models"
	"github.com/thedevrems/tuto_lua/server/internal/token"
)

// UserLoader fetches the full user record from a token's subject id.
type UserLoader interface {
	GetUserByID(id string) (models.User, error)
}

// Middleware guards routes using bearer tokens and role checks.
type Middleware struct {
	tokens *token.Manager
	users  UserLoader
}

// NewMiddleware builds the guard from the token manager and user loader.
func NewMiddleware(tokens *token.Manager, users UserLoader) *Middleware {
	return &Middleware{tokens: tokens, users: users}
}

// RequireAuth rejects unauthenticated requests and injects the user otherwise.
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := m.authenticate(r)
		if !ok {
			httpx.Error(w, http.StatusUnauthorized, "authentification requise")
			return
		}
		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
	})
}

// RequireAdmin allows the request only for authenticated admins.
func (m *Middleware) RequireAdmin(next http.Handler) http.Handler {
	return m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, _ := UserFrom(r.Context()); !user.IsAdmin() {
			httpx.Error(w, http.StatusForbidden, "accès réservé aux administrateurs")
			return
		}
		next.ServeHTTP(w, r)
	}))
}

// authenticate validates the bearer token and loads the matching user.
func (m *Middleware) authenticate(r *http.Request) (models.User, bool) {
	raw := bearerToken(r)
	if raw == "" {
		return models.User{}, false
	}
	claims, err := m.tokens.Parse(raw)
	if err != nil {
		return models.User{}, false
	}
	user, err := m.users.GetUserByID(claims.UserID)
	if err != nil {
		return models.User{}, false
	}
	return user, true
}

// bearerToken extracts the token from an "Authorization: Bearer <t>" header.
func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if len(h) > len(prefix) && strings.EqualFold(h[:len(prefix)], prefix) {
		return strings.TrimSpace(h[len(prefix):])
	}
	return ""
}
