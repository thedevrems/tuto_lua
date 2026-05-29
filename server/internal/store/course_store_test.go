package store

import (
	"errors"
	"testing"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

// seedTree inserts a small course → chapter → lesson + exercise + test tree.
func seedTree(t *testing.T, s *Store, slug string, published bool) {
	t.Helper()
	courseID, err := s.CreateCourse(models.Course{Slug: slug, Title: "T", Currency: "eur", Published: published})
	if err != nil {
		t.Fatalf("CreateCourse: %v", err)
	}
	chID, err := s.CreateChapter(models.Chapter{CourseID: courseID, Title: "C1"})
	if err != nil {
		t.Fatalf("CreateChapter: %v", err)
	}
	if _, err := s.CreateLesson(models.Lesson{ChapterID: chID, Title: "Cours", Content: "# Hi"}); err != nil {
		t.Fatalf("CreateLesson: %v", err)
	}
	exID, err := s.CreateExercise(models.Exercise{ChapterID: chID, Title: "Ex1", Difficulty: "facile", Solution: "print(1)"})
	if err != nil {
		t.Fatalf("CreateExercise: %v", err)
	}
	if _, err := s.CreateTest(models.ExerciseTest{ExerciseID: exID, Name: "t1", Code: "assert(true)"}); err != nil {
		t.Fatalf("CreateTest: %v", err)
	}
}

func TestGetCourseTreeLoadsNestedContent(t *testing.T) {
	s := newTestStore(t)
	seedTree(t, s, "m1", true)

	course, err := s.GetCourseTree("m1")
	if err != nil {
		t.Fatalf("GetCourseTree: %v", err)
	}
	if len(course.Chapters) != 1 {
		t.Fatalf("chapters = %d, want 1", len(course.Chapters))
	}
	ch := course.Chapters[0]
	if len(ch.Lessons) != 1 || ch.Lessons[0].Content != "# Hi" {
		t.Fatalf("lessons = %+v", ch.Lessons)
	}
	if len(ch.Exercises) != 1 || len(ch.Exercises[0].Tests) != 1 {
		t.Fatalf("exercises/tests not loaded: %+v", ch.Exercises)
	}
	if ch.Exercises[0].Solution != "print(1)" {
		t.Fatalf("solution = %q", ch.Exercises[0].Solution)
	}
}

func TestGetCourseTreeNotFound(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetCourseTree("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListCoursesPublishedOnly(t *testing.T) {
	s := newTestStore(t)
	seedTree(t, s, "pub", true)
	seedTree(t, s, "draft", false)

	all, _ := s.ListCourses(false)
	if len(all) != 2 {
		t.Fatalf("all courses = %d, want 2", len(all))
	}
	published, _ := s.ListCourses(true)
	if len(published) != 1 || published[0].Slug != "pub" {
		t.Fatalf("published = %+v, want only 'pub'", published)
	}
}
