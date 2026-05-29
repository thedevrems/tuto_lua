package auth

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/thedevrems/tuto_lua/server/internal/models"
	"github.com/thedevrems/tuto_lua/server/internal/store"
	"github.com/thedevrems/tuto_lua/server/internal/token"
)

// fakeUserStore is an in-memory UserStore for fast, isolated service tests.
type fakeUserStore struct {
	users map[string]models.User // keyed by id
}

func newFakeStore() *fakeUserStore { return &fakeUserStore{users: map[string]models.User{}} }

func (f *fakeUserStore) CreateUser(username, email, hash string, role models.Role) (models.User, error) {
	for _, u := range f.users {
		if strings.EqualFold(u.Username, username) || strings.EqualFold(u.Email, email) {
			return models.User{}, store.ErrConflict
		}
	}
	u := models.User{ID: username, Username: username, Email: email, PasswordHash: hash, Role: role}
	f.users[u.ID] = u
	return u, nil
}

func (f *fakeUserStore) GetUserByEmail(email string) (models.User, error) {
	for _, u := range f.users {
		if strings.EqualFold(u.Email, email) {
			return u, nil
		}
	}
	return models.User{}, store.ErrNotFound
}

func (f *fakeUserStore) GetUserByUsername(name string) (models.User, error) {
	for _, u := range f.users {
		if strings.EqualFold(u.Username, name) {
			return u, nil
		}
	}
	return models.User{}, store.ErrNotFound
}

func (f *fakeUserStore) CountUsers() (int, error) { return len(f.users), nil }

func newTestService() *Service {
	return NewService(newFakeStore(), token.NewManager("secret", time.Hour))
}

func TestRegisterFirstUserIsAdmin(t *testing.T) {
	svc := newTestService()
	res, err := svc.Register("alice", "alice@example.com", "Password1")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if res.User.Role != models.RoleAdmin {
		t.Fatalf("first user role = %q, want admin", res.User.Role)
	}
	if res.Token == "" {
		t.Fatal("expected a token")
	}
}

func TestRegisterSecondUserIsPlain(t *testing.T) {
	svc := newTestService()
	_, _ = svc.Register("alice", "alice@example.com", "Password1")
	res, err := svc.Register("bob", "bob@example.com", "Password1")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if res.User.Role != models.RoleUser {
		t.Fatalf("second user role = %q, want user", res.User.Role)
	}
}

func TestRegisterRejectsBadInput(t *testing.T) {
	svc := newTestService()
	if _, err := svc.Register("ab", "alice@example.com", "Password1"); err == nil {
		t.Fatal("expected username validation error")
	}
	if _, err := svc.Register("alice", "not-an-email", "Password1"); err == nil {
		t.Fatal("expected email validation error")
	}
	if _, err := svc.Register("alice", "alice@example.com", "weak"); err == nil {
		t.Fatal("expected password validation error")
	}
}

func TestRegisterDuplicateConflicts(t *testing.T) {
	svc := newTestService()
	_, _ = svc.Register("alice", "alice@example.com", "Password1")
	_, err := svc.Register("alice", "other@example.com", "Password1")
	if !errors.Is(err, store.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestLoginSucceedsByEmailAndUsername(t *testing.T) {
	svc := newTestService()
	_, _ = svc.Register("alice", "alice@example.com", "Password1")

	if _, err := svc.Login("alice@example.com", "Password1"); err != nil {
		t.Fatalf("login by email: %v", err)
	}
	if _, err := svc.Login("alice", "Password1"); err != nil {
		t.Fatalf("login by username: %v", err)
	}
}

func TestLoginRejectsWrongPasswordAndUnknownUser(t *testing.T) {
	svc := newTestService()
	_, _ = svc.Register("alice", "alice@example.com", "Password1")

	if _, err := svc.Login("alice@example.com", "WrongPass1"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("wrong password: got %v", err)
	}
	if _, err := svc.Login("ghost@example.com", "Password1"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("unknown user: got %v", err)
	}
}
