package token

import (
	"errors"
	"testing"
	"time"
)

func TestIssueThenParseRoundTrips(t *testing.T) {
	m := NewManager("test-secret", time.Hour)
	raw, err := m.Issue("user-1", "admin")
	if err != nil {
		t.Fatalf("Issue error: %v", err)
	}
	claims, err := m.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if claims.UserID != "user-1" || claims.Role != "admin" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	issuer := NewManager("secret-a", time.Hour)
	raw, _ := issuer.Issue("user-1", "user")

	verifier := NewManager("secret-b", time.Hour)
	if _, err := verifier.Parse(raw); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid for wrong secret, got %v", err)
	}
}

func TestParseRejectsExpiredToken(t *testing.T) {
	m := NewManager("test-secret", -time.Minute) // already expired
	raw, _ := m.Issue("user-1", "user")
	if _, err := m.Parse(raw); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid for expired token, got %v", err)
	}
}

func TestParseRejectsGarbage(t *testing.T) {
	m := NewManager("test-secret", time.Hour)
	if _, err := m.Parse("not.a.jwt"); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid for garbage, got %v", err)
	}
}
