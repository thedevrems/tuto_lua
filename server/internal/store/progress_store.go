package store

import (
	"database/sql"
	"errors"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

const progressColumns = `id, user_id, exercise_id, code, completed, updated_at`

// UpsertProgress saves a user's latest code for an exercise. The completed flag
// is monotonic: once an exercise is solved it stays solved.
func (s *Store) UpsertProgress(userID, exerciseID, code string, completed bool) (models.Progress, error) {
	_, err := s.db.Exec(
		`INSERT INTO progress (id, user_id, exercise_id, code, completed)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(user_id, exercise_id) DO UPDATE SET
		   code = excluded.code,
		   completed = MAX(progress.completed, excluded.completed),
		   updated_at = CURRENT_TIMESTAMP`,
		newID(), userID, exerciseID, code, boolToInt(completed))
	if err != nil {
		return models.Progress{}, err
	}
	return s.GetProgress(userID, exerciseID)
}

// GetProgress returns one user's state for a single exercise.
func (s *Store) GetProgress(userID, exerciseID string) (models.Progress, error) {
	row := s.db.QueryRow(
		`SELECT `+progressColumns+` FROM progress WHERE user_id = ? AND exercise_id = ?`, userID, exerciseID)
	return scanProgress(row)
}

// ListProgress returns every progress row for a user (their whole history).
func (s *Store) ListProgress(userID string) ([]models.Progress, error) {
	rows, err := s.db.Query(`SELECT `+progressColumns+` FROM progress WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []models.Progress
	for rows.Next() {
		p, err := scanProgress(rows)
		if err != nil {
			return nil, err
		}
		all = append(all, p)
	}
	return all, rows.Err()
}

// scanProgress reads one progress row, mapping the 0/1 completed flag to a bool.
func scanProgress(row rowScanner) (models.Progress, error) {
	var p models.Progress
	var completed int
	err := row.Scan(&p.ID, &p.UserID, &p.ExerciseID, &p.Code, &completed, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Progress{}, ErrNotFound
	}
	if err != nil {
		return models.Progress{}, err
	}
	p.Completed = completed != 0
	return p, nil
}
