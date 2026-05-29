// Package validate centralises input rules for registration and content.
package validate

import (
	"errors"
	"regexp"
	"strings"
)

// Validation errors are returned verbatim to the client, so keep them clear.
var (
	ErrUsername = errors.New("le nom d'utilisateur doit faire 3 à 30 caractères (lettres, chiffres, _ ou -)")
	ErrEmail    = errors.New("adresse e-mail invalide")
	ErrPassword = errors.New("le mot de passe doit faire au moins 8 caractères, avec une majuscule, une minuscule et un chiffre")
)

var (
	usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`)
	emailRe    = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

// Username trims and checks the handle, returning the normalized value.
func Username(raw string) (string, error) {
	u := strings.TrimSpace(raw)
	if !usernameRe.MatchString(u) {
		return "", ErrUsername
	}
	return u, nil
}

// Email lowercases, trims and validates the address.
func Email(raw string) (string, error) {
	e := strings.ToLower(strings.TrimSpace(raw))
	if len(e) > 254 || !emailRe.MatchString(e) {
		return "", ErrEmail
	}
	return e, nil
}

// Password enforces a minimum strength without mutating the value.
func Password(raw string) error {
	if len(raw) < 8 || len(raw) > 128 {
		return ErrPassword
	}
	if !hasUpper(raw) || !hasLower(raw) || !hasDigit(raw) {
		return ErrPassword
	}
	return nil
}

func hasUpper(s string) bool { return strings.IndexFunc(s, isUpper) >= 0 }
func hasLower(s string) bool { return strings.IndexFunc(s, isLower) >= 0 }
func hasDigit(s string) bool { return strings.IndexFunc(s, isDigit) >= 0 }

func isUpper(r rune) bool { return r >= 'A' && r <= 'Z' }
func isLower(r rune) bool { return r >= 'a' && r <= 'z' }
func isDigit(r rune) bool { return r >= '0' && r <= '9' }
