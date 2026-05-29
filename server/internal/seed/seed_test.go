package seed_test

import (
	"testing"

	"github.com/thedevrems/tuto_lua/server/internal/database"
	"github.com/thedevrems/tuto_lua/server/internal/seed"
	"github.com/thedevrems/tuto_lua/server/internal/store"
)

func newStore(t *testing.T) *store.Store {
	t.Helper()
	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return store.New(db)
}

func TestRunPopulatesCatalogue(t *testing.T) {
	s := newStore(t)
	if err := seed.Run(s); err != nil {
		t.Fatalf("seed.Run: %v", err)
	}
	courses, err := s.ListCourses(true)
	if err != nil {
		t.Fatalf("ListCourses: %v", err)
	}
	if len(courses) == 0 {
		t.Fatal("expected the seed to insert at least one course")
	}

	// The first course should load as a full tree with content.
	tree, err := s.GetCourseTree(courses[0].Slug)
	if err != nil {
		t.Fatalf("GetCourseTree: %v", err)
	}
	if len(tree.Chapters) == 0 {
		t.Fatal("expected the first course to have chapters")
	}
}

func TestRunIsIdempotent(t *testing.T) {
	s := newStore(t)
	if err := seed.Run(s); err != nil {
		t.Fatalf("first run: %v", err)
	}
	first, _ := s.CountCourses()
	if err := seed.Run(s); err != nil {
		t.Fatalf("second run: %v", err)
	}
	second, _ := s.CountCourses()
	if first != second {
		t.Fatalf("seed not idempotent: %d then %d courses", first, second)
	}
}
