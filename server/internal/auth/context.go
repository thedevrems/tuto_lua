// Package auth orchestrates registration, login and request authorization.
package auth

import (
	"context"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

// ctxKey is unexported so only this package can write the user into a context.
type ctxKey struct{}

var userCtxKey = ctxKey{}

// WithUser returns a copy of ctx carrying the authenticated user.
func WithUser(ctx context.Context, u models.User) context.Context {
	return context.WithValue(ctx, userCtxKey, u)
}

// UserFrom extracts the authenticated user previously stored by the middleware.
func UserFrom(ctx context.Context) (models.User, bool) {
	u, ok := ctx.Value(userCtxKey).(models.User)
	return u, ok
}
