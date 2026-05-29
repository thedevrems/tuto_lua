package store

import (
	"errors"
	"testing"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

func TestCreateAndFetchUser(t *testing.T) {
	s := newTestStore(t)
	created, err := s.CreateUser("alice", "alice@example.com", "hash", models.RoleUser)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if created.ID == "" || created.Role != models.RoleUser {
		t.Fatalf("unexpected created user: %+v", created)
	}

	byEmail, err := s.GetUserByEmail("ALICE@example.com") // case-insensitive
	if err != nil || byEmail.ID != created.ID {
		t.Fatalf("GetUserByEmail = (%+v, %v)", byEmail, err)
	}
	byName, err := s.GetUserByUsername("Alice")
	if err != nil || byName.ID != created.ID {
		t.Fatalf("GetUserByUsername = (%+v, %v)", byName, err)
	}
}

func TestCreateUserDuplicateEmailConflicts(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.CreateUser("alice", "dup@example.com", "h", models.RoleUser); err != nil {
		t.Fatalf("first insert: %v", err)
	}
	_, err := s.CreateUser("bob", "DUP@example.com", "h", models.RoleUser)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict on duplicate email, got %v", err)
	}
}

func TestCreateUserDuplicateUsernameConflicts(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.CreateUser("alice", "a@example.com", "h", models.RoleUser); err != nil {
		t.Fatalf("first insert: %v", err)
	}
	_, err := s.CreateUser("Alice", "b@example.com", "h", models.RoleUser)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict on duplicate username, got %v", err)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetUserByID("does-not-exist"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCountUsers(t *testing.T) {
	s := newTestStore(t)
	if n, _ := s.CountUsers(); n != 0 {
		t.Fatalf("empty store count = %d, want 0", n)
	}
	_, _ = s.CreateUser("alice", "a@example.com", "h", models.RoleUser)
	if n, _ := s.CountUsers(); n != 1 {
		t.Fatalf("count after insert = %d, want 1", n)
	}
}
