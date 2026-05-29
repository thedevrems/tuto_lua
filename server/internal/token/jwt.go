// Package token issues and validates the JWTs used for stateless auth.
package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalid is returned for any token that fails signature or claim checks.
var ErrInvalid = errors.New("invalid token")

// Claims is the authenticated identity carried by a token.
type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
}

// Manager signs and verifies tokens with a single HMAC secret.
type Manager struct {
	secret []byte
	ttl    time.Duration
}

// NewManager builds a Manager from the configured secret and lifetime.
func NewManager(secret string, ttl time.Duration) *Manager {
	return &Manager{secret: []byte(secret), ttl: ttl}
}

// Issue returns a signed token embedding the user id and role.
func (m *Manager) Issue(userID, role string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"uid":  userID,
		"role": role,
		"iat":  now.Unix(),
		"exp":  now.Add(m.ttl).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(m.secret)
}

// Parse validates the token signature and expiry, returning its claims.
func (m *Manager) Parse(raw string) (Claims, error) {
	parsed, err := jwt.Parse(raw, m.keyFunc, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !parsed.Valid {
		return Claims{}, ErrInvalid
	}
	return claimsFrom(parsed)
}

// keyFunc supplies the HMAC secret to the parser.
func (m *Manager) keyFunc(*jwt.Token) (any, error) {
	return m.secret, nil
}

// claimsFrom extracts our typed Claims from the parsed token.
func claimsFrom(parsed *jwt.Token) (Claims, error) {
	mc, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, ErrInvalid
	}
	uid, _ := mc["uid"].(string)
	role, _ := mc["role"].(string)
	if uid == "" {
		return Claims{}, ErrInvalid
	}
	return Claims{UserID: uid, Role: role}, nil
}
