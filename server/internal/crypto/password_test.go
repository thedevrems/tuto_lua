package crypto

import (
	"errors"
	"testing"
)

func TestHashPasswordProducesVerifiableHash(t *testing.T) {
	hash, err := HashPassword("Sup3rSecret!")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "Sup3rSecret!" {
		t.Fatal("hash must not equal the plaintext password")
	}
	if err := VerifyPassword(hash, "Sup3rSecret!"); err != nil {
		t.Fatalf("VerifyPassword rejected the correct password: %v", err)
	}
}

func TestVerifyPasswordRejectsWrongPassword(t *testing.T) {
	hash, err := HashPassword("correct-horse")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if err := VerifyPassword(hash, "wrong-horse"); !errors.Is(err, ErrMismatch) {
		t.Fatalf("expected ErrMismatch, got %v", err)
	}
}

func TestHashPasswordIsSaltedPerCall(t *testing.T) {
	first, _ := HashPassword("same-input")
	second, _ := HashPassword("same-input")
	if first == second {
		t.Fatal("two hashes of the same password should differ (random salt)")
	}
}
