package models

import "time"

// Notification is an in-app message for a single user.
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Link      string    `json:"link"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"createdAt"`
}
