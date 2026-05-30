package store

import "github.com/thedevrems/tuto_lua/server/internal/models"

// This file holds the admin update/delete operations for course content.
// Deletes cascade to children via the schema's ON DELETE CASCADE rules.

// UpdateCourse edits a course's fields (by id). A duplicate slug maps to ErrConflict.
func (s *Store) UpdateCourse(c models.Course) error {
	err := s.execAffecting(
		`UPDATE courses SET slug=?, title=?, summary=?, price_cents=?, currency=?, published=?, position=? WHERE id=?`,
		c.Slug, c.Title, c.Summary, c.PriceCents, c.Currency, boolToInt(c.Published), c.Position, c.ID)
	if isUniqueViolation(err) {
		return ErrConflict
	}
	return err
}

// DeleteCourse removes a course and (via cascade) all its content.
func (s *Store) DeleteCourse(id string) error {
	return s.execAffecting(`DELETE FROM courses WHERE id=?`, id)
}

// UpdateChapter edits a chapter's title, summary and position.
func (s *Store) UpdateChapter(c models.Chapter) error {
	return s.execAffecting(`UPDATE chapters SET title=?, summary=?, position=? WHERE id=?`, c.Title, c.Summary, c.Position, c.ID)
}

// DeleteChapter removes a chapter and its lessons/exercises.
func (s *Store) DeleteChapter(id string) error {
	return s.execAffecting(`DELETE FROM chapters WHERE id=?`, id)
}

// UpdateLesson edits a lesson's title, content and position.
func (s *Store) UpdateLesson(l models.Lesson) error {
	return s.execAffecting(`UPDATE lessons SET title=?, content=?, position=? WHERE id=?`, l.Title, l.Content, l.Position, l.ID)
}

// DeleteLesson removes a lesson.
func (s *Store) DeleteLesson(id string) error {
	return s.execAffecting(`DELETE FROM lessons WHERE id=?`, id)
}

// UpdateExercise edits an exercise (hints are stored as JSON).
func (s *Store) UpdateExercise(e models.Exercise) error {
	return s.execAffecting(
		`UPDATE exercises SET title=?, difficulty=?, statement=?, starter=?, solution=?, hints=?, position=? WHERE id=?`,
		e.Title, e.Difficulty, e.Statement, e.Starter, e.Solution, encodeHints(e.Hints), e.Position, e.ID)
}

// DeleteExercise removes an exercise and its tests.
func (s *Store) DeleteExercise(id string) error {
	return s.execAffecting(`DELETE FROM exercises WHERE id=?`, id)
}

// UpdateTest edits an automated test's name, code and position.
func (s *Store) UpdateTest(t models.ExerciseTest) error {
	return s.execAffecting(`UPDATE exercise_tests SET name=?, code=?, position=? WHERE id=?`, t.Name, t.Code, t.Position, t.ID)
}

// DeleteTest removes an automated test.
func (s *Store) DeleteTest(id string) error {
	return s.execAffecting(`DELETE FROM exercise_tests WHERE id=?`, id)
}
