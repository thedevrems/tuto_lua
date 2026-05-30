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
	Courses       *handlers.CourseHandler
	Progress      *handlers.ProgressHandler
	Enrollments   *handlers.EnrollmentHandler
	Admin         *handlers.AdminHandler
	Payments      *handlers.PaymentHandler
	Profile       *handlers.ProfileHandler
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
		mountCourseRoutes(api, d)
		mountProgressRoutes(api, d)
		mountProfileRoutes(api, d)
		mountAdminRoutes(api, d)
		mountPaymentRoutes(api, d)
	})
	return r
}

// mountProfileRoutes exposes the authenticated user's own account endpoints.
func mountProfileRoutes(api chi.Router, d Deps) {
	api.Route("/me", func(m chi.Router) {
		m.Use(d.Guard.RequireAuth)
		m.Get("/courses", d.Profile.MyCourses)
		m.Post("/password", d.Profile.ChangePassword)
	})
}

// mountPaymentRoutes exposes Stripe Checkout creation (authenticated) and the
// Stripe webhook (public — verified by signature, not by a token).
func mountPaymentRoutes(api chi.Router, d Deps) {
	api.Route("/payments", func(p chi.Router) {
		p.With(d.Guard.RequireAuth).Post("/checkout", d.Payments.Checkout)
		p.Post("/webhook", d.Payments.Webhook)
	})
}

// mountAdminRoutes exposes the admin-only management & authoring endpoints.
func mountAdminRoutes(api chi.Router, d Deps) {
	api.Route("/admin", func(a chi.Router) {
		a.Use(d.Guard.RequireAdmin)
		a.Get("/users", d.Admin.ListUsers)
		a.Get("/users/{userId}/progress", d.Admin.UserProgress)
		a.Get("/users/{userId}/courses", d.Admin.UserCourses)
		a.Post("/enrollments", d.Admin.GrantAccess)
		a.Get("/courses", d.Admin.ListCourses)
		a.Post("/courses", d.Admin.CreateCourse)
		a.Post("/courses/{courseId}/chapters", d.Admin.CreateChapter)
		a.Post("/chapters/{chapterId}/lessons", d.Admin.CreateLesson)
		a.Post("/chapters/{chapterId}/exercises", d.Admin.CreateExercise)
		a.Post("/exercises/{exerciseId}/tests", d.Admin.CreateTest)
		mountAdminContentEdits(a, d)
	})
}

// mountAdminContentEdits adds the update (PUT) and delete (DELETE) routes for
// every content type. Deletes cascade to children in the database.
func mountAdminContentEdits(a chi.Router, d Deps) {
	a.Put("/courses/{courseId}", d.Admin.UpdateCourse)
	a.Delete("/courses/{courseId}", d.Admin.DeleteCourse)
	a.Put("/chapters/{chapterId}", d.Admin.UpdateChapter)
	a.Delete("/chapters/{chapterId}", d.Admin.DeleteChapter)
	a.Put("/lessons/{lessonId}", d.Admin.UpdateLesson)
	a.Delete("/lessons/{lessonId}", d.Admin.DeleteLesson)
	a.Put("/exercises/{exerciseId}", d.Admin.UpdateExercise)
	a.Delete("/exercises/{exerciseId}", d.Admin.DeleteExercise)
	a.Put("/tests/{testId}", d.Admin.UpdateTest)
	a.Delete("/tests/{testId}", d.Admin.DeleteTest)
}

// mountCourseRoutes exposes the public course catalogue and content tree.
// The tree uses optional auth so paid content can be gated per user.
func mountCourseRoutes(api chi.Router, d Deps) {
	api.Route("/courses", func(c chi.Router) {
		c.Get("/", d.Courses.List)
		c.With(d.Guard.OptionalAuth).Get("/{slug}", d.Courses.Tree)
	})
}

// mountProgressRoutes exposes the per-user progress and enrollment endpoints.
func mountProgressRoutes(api chi.Router, d Deps) {
	api.Route("/progress", func(p chi.Router) {
		p.Use(d.Guard.RequireAuth)
		p.Get("/", d.Progress.List)
		p.Put("/{exerciseId}", d.Progress.Save)
	})
	api.With(d.Guard.RequireAuth).Get("/enrollments", d.Enrollments.Mine)
}

// mountAuthRoutes groups the public and protected authentication endpoints.
func mountAuthRoutes(api chi.Router, d Deps) {
	api.Route("/auth", func(a chi.Router) {
		a.Post("/register", d.Auth.Register)
		a.Post("/login", d.Auth.Login)
		a.With(d.Guard.RequireAuth).Get("/me", d.Auth.Me)
	})
}
