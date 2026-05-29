// Package router assembles the HTTP routes and shared middleware.
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/thedevrems/tuto_lua/server/internal/auth"
	"github.com/thedevrems/tuto_lua/server/internal/handlers"
)

// Deps are the wired-up collaborators the routes need.
type Deps struct {
	AllowedOrigin string
	Auth          *handlers.AuthHandler
	Guard         *auth.Middleware
}

// New builds the application's HTTP handler with all routes mounted under /api.
func New(d Deps) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)
	r.Use(cors(d.AllowedOrigin))

	r.Route("/api", func(api chi.Router) {
		api.Get("/health", handlers.Health)
		mountAuthRoutes(api, d)
	})
	return r
}

// mountAuthRoutes groups the public and protected authentication endpoints.
func mountAuthRoutes(api chi.Router, d Deps) {
	api.Route("/auth", func(a chi.Router) {
		a.Post("/register", d.Auth.Register)
		a.Post("/login", d.Auth.Login)
		a.With(d.Guard.RequireAuth).Get("/me", d.Auth.Me)
	})
}
