package store

import (
	"errors"
	"testing"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

// makeTest builds a course→chapter→exercise→test and returns the ids needed.
func makeTest(t *testing.T, s *Store) (slug, exerciseID, testID string) {
	t.Helper()
	courseID, _ := s.CreateCourse(models.Course{Slug: "c", Title: "C", Currency: "eur"})
	chID, _ := s.CreateChapter(models.Chapter{CourseID: courseID, Title: "Ch"})
	exID, _ := s.CreateExercise(models.Exercise{ChapterID: chID, Title: "Ex"})
	tID, err := s.CreateTest(models.ExerciseTest{ExerciseID: exID, Name: "t1", Code: "assert(true)"})
	if err != nil {
		t.Fatalf("CreateTest: %v", err)
	}
	return "c", exID, tID
}

func TestUpdateTestChangesFields(t *testing.T) {
	s := newTestStore(t)
	slug, _, testID := makeTest(t, s)

	if err := s.UpdateTest(models.ExerciseTest{ID: testID, Name: "renommé", Code: "assert(1==1)", Position: 0}); err != nil {
		t.Fatalf("UpdateTest: %v", err)
	}
	tree, _ := s.GetCourseTree(slug)
	got := tree.Chapters[0].Exercises[0].Tests[0]
	if got.Name != "renommé" || got.Code != "assert(1==1)" {
		t.Fatalf("test not updated: %+v", got)
	}
}

func TestDeleteTestRemovesIt(t *testing.T) {
	s := newTestStore(t)
	slug, _, testID := makeTest(t, s)

	if err := s.DeleteTest(testID); err != nil {
		t.Fatalf("DeleteTest: %v", err)
	}
	tree, _ := s.GetCourseTree(slug)
	if n := len(tree.Chapters[0].Exercises[0].Tests); n != 0 {
		t.Fatalf("test still present after delete: %d", n)
	}
}

func TestUpdateAndDeleteMissingReturnNotFound(t *testing.T) {
	s := newTestStore(t)
	if err := s.DeleteTest("nope"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("DeleteTest(missing) = %v, want ErrNotFound", err)
	}
	if err := s.UpdateExercise(models.Exercise{ID: "nope", Title: "x"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("UpdateExercise(missing) = %v, want ErrNotFound", err)
	}
}

func TestDeleteExerciseCascadesTests(t *testing.T) {
	s := newTestStore(t)
	slug, exerciseID, _ := makeTest(t, s)

	if err := s.DeleteExercise(exerciseID); err != nil {
		t.Fatalf("DeleteExercise: %v", err)
	}
	tree, _ := s.GetCourseTree(slug)
	if n := len(tree.Chapters[0].Exercises); n != 0 {
		t.Fatalf("exercise still present after delete: %d", n)
	}
}
