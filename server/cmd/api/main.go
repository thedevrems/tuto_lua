// Command api is the Lua Academy HTTP backend entrypoint.
package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thedevrems/tuto_lua/server/internal/auth"
	"github.com/thedevrems/tuto_lua/server/internal/config"
	"github.com/thedevrems/tuto_lua/server/internal/database"
	"github.com/thedevrems/tuto_lua/server/internal/handlers"
	"github.com/thedevrems/tuto_lua/server/internal/payment"
	"github.com/thedevrems/tuto_lua/server/internal/router"
	"github.com/thedevrems/tuto_lua/server/internal/seed"
	"github.com/thedevrems/tuto_lua/server/internal/store"
	"github.com/thedevrems/tuto_lua/server/internal/token"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal: %v", err)
	}
}

// run wires every dependency and blocks until the server is shut down.
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	db, err := openDatabase(cfg.DatabasePath)
	if err != nil {
		return err
	}
	defer db.Close()

	st := store.New(db)
	if err := seed.Run(st); err != nil {
		return err
	}
	return serve(cfg.Port, buildRouter(cfg, st))
}

// openDatabase opens the connection and applies the schema once.
func openDatabase(path string) (*sql.DB, error) {
	db, err := database.Open(path)
	if err != nil {
		return nil, err
	}
	if err := database.Migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// buildRouter constructs the HTTP handler from the store and config.
func buildRouter(cfg config.Config, st *store.Store) http.Handler {
	tokens := token.NewManager(cfg.JWTSecret, cfg.JWTTTL)
	authSvc := auth.NewService(st, tokens)
	guard := auth.NewMiddleware(tokens, st)
	paySvc := payment.NewService(st, cfg.StripeSecretKey, cfg.StripeWebhook, cfg.FrontendURL)
	return router.New(router.Deps{
		AllowedOrigin: cfg.AllowedOrigin,
		Auth:          handlers.NewAuthHandler(authSvc),
		Courses:       handlers.NewCourseHandler(st),
		Progress:      handlers.NewProgressHandler(st),
		Enrollments:   handlers.NewEnrollmentHandler(st),
		Admin:         handlers.NewAdminHandler(st),
		Payments:      handlers.NewPaymentHandler(paySvc),
		Profile:       handlers.NewProfileHandler(authSvc, st),
		Guard:         guard,
	})
}

// serve starts the HTTP server and shuts it down gracefully on SIGINT/SIGTERM.
func serve(port string, handler http.Handler) error {
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		log.Printf("API listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("shutting down…")
	return srv.Shutdown(ctx)
}
