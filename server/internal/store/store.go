// Package store is the data-access layer: every SQL query lives here so the
// rest of the app depends on Go methods, not on the database.
package store

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
)

// Sentinel errors let callers branch on outcome without inspecting SQL.
var (
	ErrNotFound = errors.New("ressource introuvable")
	ErrConflict = errors.New("ressource déjà existante")
)

// Store wraps the database handle and exposes typed repository methods.
type Store struct {
	db *sql.DB
}

// New builds a Store around an open database connection.
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// newID returns a fresh random identifier for a new row.
func newID() string {
	return uuid.NewString()
}

// isUniqueViolation reports whether err is a UNIQUE constraint failure,
// which we translate into ErrConflict for the caller.
func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
