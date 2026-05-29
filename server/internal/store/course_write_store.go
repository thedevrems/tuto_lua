package store

import (
	"encoding/json"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

// CountCourses returns the number of courses (used to gate seeding).
func (s *Store) CountCourses() (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM courses`).Scan(&n)
	return n, err
}

// CreateCourse inserts a course and returns its generated id.
func (s *Store) CreateCourse(c models.Course) (string, error) {
	id := newID()
	_, err := s.db.Exec(
		`INSERT INTO courses (id, slug, title, summary, price_cents, currency, published, position)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, c.Slug, c.Title, c.Summary, c.PriceCents, c.Currency, boolToInt(c.Published), c.Position)
	if isUniqueViolation(err) {
		return "", ErrConflict
	}
	return id, err
}

// CreateChapter inserts a chapter under a course.
func (s *Store) CreateChapter(c models.Chapter) (string, error) {
	id := newID()
	_, err := s.db.Exec(
		`INSERT INTO chapters (id, course_id, title, summary, position) VALUES (?, ?, ?, ?, ?)`,
		id, c.CourseID, c.Title, c.Summary, c.Position)
	return id, err
}

// CreateLesson inserts a lesson under a chapter.
func (s *Store) CreateLesson(l models.Lesson) (string, error) {
	id := newID()
	_, err := s.db.Exec(
		`INSERT INTO lessons (id, chapter_id, title, content, position) VALUES (?, ?, ?, ?, ?)`,
		id, l.ChapterID, l.Title, l.Content, l.Position)
	return id, err
}

// CreateExercise inserts an exercise under a chapter, storing hints as JSON.
func (s *Store) CreateExercise(e models.Exercise) (string, error) {
	id := newID()
	_, err := s.db.Exec(
		`INSERT INTO exercises (id, chapter_id, title, difficulty, statement, starter, solution, hints, position)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, e.ChapterID, e.Title, e.Difficulty, e.Statement, e.Starter, e.Solution, encodeHints(e.Hints), e.Position)
	return id, err
}

// encodeHints serialises hints to a JSON array, never failing (defaults to []).
func encodeHints(hints []string) string {
	if len(hints) == 0 {
		return "[]"
	}
	raw, err := json.Marshal(hints)
	if err != nil {
		return "[]"
	}
	return string(raw)
}

// CreateTest inserts an automated test under an exercise.
func (s *Store) CreateTest(t models.ExerciseTest) (string, error) {
	id := newID()
	_, err := s.db.Exec(
		`INSERT INTO exercise_tests (id, exercise_id, name, code, position) VALUES (?, ?, ?, ?, ?)`,
		id, t.ExerciseID, t.Name, t.Code, t.Position)
	return id, err
}

// boolToInt maps a Go bool to SQLite's 0/1 integer convention.
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
