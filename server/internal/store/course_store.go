package store

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

// rowScanner is satisfied by both *sql.Row and *sql.Rows.
type rowScanner interface{ Scan(dest ...any) error }

const courseColumns = `id, slug, title, summary, price_cents, currency, published, position, created_at`

// ListCourses returns courses ordered by position; published-only when asked.
func (s *Store) ListCourses(publishedOnly bool) ([]models.Course, error) {
	query := `SELECT ` + courseColumns + ` FROM courses`
	if publishedOnly {
		query += ` WHERE published = 1`
	}
	query += ` ORDER BY position, created_at`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		c, err := scanCourse(rows)
		if err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, rows.Err()
}

// GetCourseBySlug loads a single course header (no nested content).
func (s *Store) GetCourseBySlug(slug string) (models.Course, error) {
	return scanCourse(s.db.QueryRow(`SELECT `+courseColumns+` FROM courses WHERE slug = ?`, slug))
}

// GetCourseTree loads a course with all its chapters, lessons, exercises and tests.
func (s *Store) GetCourseTree(slug string) (models.Course, error) {
	course, err := s.GetCourseBySlug(slug)
	if err != nil {
		return models.Course{}, err
	}
	course.Chapters, err = s.chaptersWithContent(course.ID)
	return course, err
}

// chaptersWithContent loads every chapter of a course, fully populated.
// Chapters are read first (closing the result set) before issuing the nested
// lesson/exercise queries — required because the pool holds a single connection.
func (s *Store) chaptersWithContent(courseID string) ([]models.Chapter, error) {
	chapters, err := s.scanChapters(courseID)
	if err != nil {
		return nil, err
	}
	for i := range chapters {
		if err := s.fillChapter(&chapters[i]); err != nil {
			return nil, err
		}
	}
	return chapters, nil
}

// scanChapters reads the chapter rows of a course into memory and closes them.
func (s *Store) scanChapters(courseID string) ([]models.Chapter, error) {
	rows, err := s.db.Query(
		`SELECT id, course_id, title, summary, position FROM chapters WHERE course_id = ? ORDER BY position`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []models.Chapter
	for rows.Next() {
		var c models.Chapter
		if err := rows.Scan(&c.ID, &c.CourseID, &c.Title, &c.Summary, &c.Position); err != nil {
			return nil, err
		}
		chapters = append(chapters, c)
	}
	return chapters, rows.Err()
}

// fillChapter loads the lessons and exercises belonging to one chapter.
func (s *Store) fillChapter(c *models.Chapter) error {
	var err error
	if c.Lessons, err = s.listLessons(c.ID); err != nil {
		return err
	}
	c.Exercises, err = s.listExercises(c.ID)
	return err
}

func (s *Store) listLessons(chapterID string) ([]models.Lesson, error) {
	rows, err := s.db.Query(
		`SELECT id, chapter_id, title, content, position FROM lessons WHERE chapter_id = ? ORDER BY position`, chapterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		var l models.Lesson
		if err := rows.Scan(&l.ID, &l.ChapterID, &l.Title, &l.Content, &l.Position); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, rows.Err()
}

// listExercises returns a chapter's exercises with their tests. Exercises are
// read first (result set closed) before the per-exercise test queries run.
func (s *Store) listExercises(chapterID string) ([]models.Exercise, error) {
	exercises, err := s.scanExercises(chapterID)
	if err != nil {
		return nil, err
	}
	for i := range exercises {
		tests, err := s.listTests(exercises[i].ID)
		if err != nil {
			return nil, err
		}
		exercises[i].Tests = tests
	}
	return exercises, nil
}

// scanExercises reads the exercise rows of a chapter into memory and closes them.
func (s *Store) scanExercises(chapterID string) ([]models.Exercise, error) {
	rows, err := s.db.Query(
		`SELECT id, chapter_id, title, difficulty, statement, starter, solution, hints, position
		 FROM exercises WHERE chapter_id = ? ORDER BY position`, chapterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var e models.Exercise
		var hints string
		if err := rows.Scan(&e.ID, &e.ChapterID, &e.Title, &e.Difficulty, &e.Statement, &e.Starter, &e.Solution, &hints, &e.Position); err != nil {
			return nil, err
		}
		e.Hints = decodeHints(hints)
		exercises = append(exercises, e)
	}
	return exercises, rows.Err()
}

// decodeHints parses the JSON-encoded hints column into a slice (nil on error).
func decodeHints(raw string) []string {
	if raw == "" {
		return nil
	}
	var hints []string
	if json.Unmarshal([]byte(raw), &hints) != nil {
		return nil
	}
	return hints
}

func (s *Store) listTests(exerciseID string) ([]models.ExerciseTest, error) {
	rows, err := s.db.Query(
		`SELECT id, exercise_id, name, code, position FROM exercise_tests WHERE exercise_id = ? ORDER BY position`, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []models.ExerciseTest
	for rows.Next() {
		var t models.ExerciseTest
		if err := rows.Scan(&t.ID, &t.ExerciseID, &t.Name, &t.Code, &t.Position); err != nil {
			return nil, err
		}
		tests = append(tests, t)
	}
	return tests, rows.Err()
}

// scanCourse reads one course row, converting the 0/1 published flag to a bool.
func scanCourse(row rowScanner) (models.Course, error) {
	var c models.Course
	var published int
	err := row.Scan(&c.ID, &c.Slug, &c.Title, &c.Summary, &c.PriceCents, &c.Currency, &published, &c.Position, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Course{}, ErrNotFound
	}
	if err != nil {
		return models.Course{}, err
	}
	c.Published = published != 0
	return c, nil
}
