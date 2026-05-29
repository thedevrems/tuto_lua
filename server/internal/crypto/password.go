// Package crypto wraps password hashing so the rest of the app never touches
// bcrypt directly.
package crypto

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// ErrMismatch is returned when a password does not match its hash.
var ErrMismatch = errors.New("password does not match")

// cost is the bcrypt work factor. 12 is a sensible 2020s default.
const cost = 12

// HashPassword turns a plaintext password into a salted bcrypt hash.
func HashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword reports whether plain matches the stored bcrypt hash.
// It returns ErrMismatch on a wrong password and other errors on bad input.
func VerifyPassword(hash, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrMismatch
	}
	return err
}
