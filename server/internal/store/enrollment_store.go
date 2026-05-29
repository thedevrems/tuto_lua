package store

import (
	"database/sql"
	"errors"
)

// IsEnrolled reports whether a user already has access to a course.
func (s *Store) IsEnrolled(userID, courseID string) (bool, error) {
	var one int
	err := s.db.QueryRow(
		`SELECT 1 FROM enrollments WHERE user_id = ? AND course_id = ? LIMIT 1`, userID, courseID).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ListEnrollmentCourseIDs returns the ids of every course a user can access.
func (s *Store) ListEnrollmentCourseIDs(userID string) ([]string, error) {
	rows, err := s.db.Query(`SELECT course_id FROM enrollments WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// CreateEnrollment grants course access, idempotently (no error if it exists).
// source is "purchase" or "admin"; grantedBy is the admin id when applicable.
func (s *Store) CreateEnrollment(userID, courseID, source string, grantedBy *string) error {
	_, err := s.db.Exec(
		`INSERT INTO enrollments (id, user_id, course_id, source, granted_by)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(user_id, course_id) DO NOTHING`,
		newID(), userID, courseID, source, grantedBy)
	return err
}
