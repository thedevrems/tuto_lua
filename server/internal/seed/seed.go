// Package seed loads the initial course catalogue into an empty database.
package seed

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

//go:embed curriculum.json
var curriculumJSON []byte

// Store is the subset of the data layer the seeder writes through.
type Store interface {
	CountCourses() (int, error)
	CreateCourse(models.Course) (string, error)
	CreateChapter(models.Chapter) (string, error)
	CreateLesson(models.Lesson) (string, error)
	CreateExercise(models.Exercise) (string, error)
	CreateTest(models.ExerciseTest) (string, error)
}

// JSON shapes mirror the output of web/scripts/export-curriculum.mjs.
type seedCourse struct {
	Slug       string        `json:"slug"`
	Title      string        `json:"title"`
	Summary    string        `json:"summary"`
	PriceCents int           `json:"priceCents"`
	Currency   string        `json:"currency"`
	Published  bool          `json:"published"`
	Position   int           `json:"position"`
	Chapters   []seedChapter `json:"chapters"`
}

type seedChapter struct {
	Title     string         `json:"title"`
	Summary   string         `json:"summary"`
	Position  int            `json:"position"`
	Lessons   []seedLesson   `json:"lessons"`
	Exercises []seedExercise `json:"exercises"`
}

type seedLesson struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Position int    `json:"position"`
}

type seedExercise struct {
	Title      string     `json:"title"`
	Difficulty string     `json:"difficulty"`
	Statement  string     `json:"statement"`
	Starter    string     `json:"starter"`
	Solution   string     `json:"solution"`
	Hints      []string   `json:"hints"`
	Position   int        `json:"position"`
	Tests      []seedTest `json:"tests"`
}

type seedTest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Position int    `json:"position"`
}

// Run inserts the embedded curriculum, but only when the catalogue is empty.
func Run(s Store) error {
	n, err := s.CountCourses()
	if err != nil {
		return err
	}
	if n > 0 {
		return nil // already seeded
	}
	var courses []seedCourse
	if err := json.Unmarshal(curriculumJSON, &courses); err != nil {
		return fmt.Errorf("parse seed: %w", err)
	}
	for _, c := range courses {
		if err := insertCourse(s, c); err != nil {
			return err
		}
	}
	return nil
}

// insertCourse persists a course and recurses into its chapters.
func insertCourse(s Store, c seedCourse) error {
	id, err := s.CreateCourse(models.Course{
		Slug: c.Slug, Title: c.Title, Summary: c.Summary,
		PriceCents: c.PriceCents, Currency: c.Currency, Published: c.Published, Position: c.Position,
	})
	if err != nil {
		return err
	}
	for _, ch := range c.Chapters {
		if err := insertChapter(s, id, ch); err != nil {
			return err
		}
	}
	return nil
}

// insertChapter persists a chapter and its lessons + exercises.
func insertChapter(s Store, courseID string, ch seedChapter) error {
	id, err := s.CreateChapter(models.Chapter{CourseID: courseID, Title: ch.Title, Summary: ch.Summary, Position: ch.Position})
	if err != nil {
		return err
	}
	for _, l := range ch.Lessons {
		if _, err := s.CreateLesson(models.Lesson{ChapterID: id, Title: l.Title, Content: l.Content, Position: l.Position}); err != nil {
			return err
		}
	}
	for _, e := range ch.Exercises {
		if err := insertExercise(s, id, e); err != nil {
			return err
		}
	}
	return nil
}

// insertExercise persists an exercise and its automated tests.
func insertExercise(s Store, chapterID string, e seedExercise) error {
	id, err := s.CreateExercise(models.Exercise{
		ChapterID: chapterID, Title: e.Title, Difficulty: e.Difficulty,
		Statement: e.Statement, Starter: e.Starter, Solution: e.Solution, Hints: e.Hints, Position: e.Position,
	})
	if err != nil {
		return err
	}
	for _, t := range e.Tests {
		if _, err := s.CreateTest(models.ExerciseTest{ExerciseID: id, Name: t.Name, Code: t.Code, Position: t.Position}); err != nil {
			return err
		}
	}
	return nil
}
