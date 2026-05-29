package models

import "time"

// Enrollment records that a user has access to a course, either through a
// purchase or an admin grant.
type Enrollment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	CourseID  string    `json:"courseId"`
	GrantedBy *string   `json:"grantedBy,omitempty"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
}

// Progress is a user's latest state on a single exercise: their most recent
// code ("last push") and whether all tests have passed.
type Progress struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	ExerciseID string    `json:"exerciseId"`
	Code       string    `json:"code"`
	Completed  bool      `json:"completed"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
