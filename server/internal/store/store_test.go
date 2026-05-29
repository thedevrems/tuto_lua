package store

import (
	"testing"

	"github.com/thedevrems/tuto_lua/server/internal/database"
)

// newTestStore spins up a private in-memory database with the full schema.
// Each test gets an isolated store; the DB is closed on cleanup.
func newTestStore(t *testing.T) *Store {
	t.Helper()
	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatalf("open in-memory db: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return New(db)
}
