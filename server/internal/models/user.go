// Package models holds the plain data structures shared across layers.
package models

import "time"

// Role enumerates the permission levels a user can hold.
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User is an account. PasswordHash is never serialized to clients.
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
}

// IsAdmin reports whether the user holds administrative rights.
func (u User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
