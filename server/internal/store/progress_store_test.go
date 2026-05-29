package store

import (
	"testing"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

// makeExercise creates a user and a course→chapter→exercise, returning ids.
func makeUserAndExercise(t *testing.T, s *Store) (userID, exerciseID string) {
	t.Helper()
	user, err := s.CreateUser("learner", "l@example.com", "h", models.RoleUser)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	courseID, _ := s.CreateCourse(models.Course{Slug: "c1", Title: "C", Currency: "eur"})
	chID, _ := s.CreateChapter(models.Chapter{CourseID: courseID, Title: "Ch"})
	exID, err := s.CreateExercise(models.Exercise{ChapterID: chID, Title: "Ex"})
	if err != nil {
		t.Fatalf("CreateExercise: %v", err)
	}
	return user.ID, exID
}

func TestUpsertProgressCreatesThenUpdatesCode(t *testing.T) {
	s := newTestStore(t)
	userID, exID := makeUserAndExercise(t, s)

	first, err := s.UpsertProgress(userID, exID, "print(1)", false)
	if err != nil {
		t.Fatalf("first upsert: %v", err)
	}
	if first.Code != "print(1)" || first.Completed {
		t.Fatalf("unexpected first state: %+v", first)
	}

	second, err := s.UpsertProgress(userID, exID, "print(2)", false)
	if err != nil {
		t.Fatalf("second upsert: %v", err)
	}
	if second.Code != "print(2)" {
		t.Fatalf("code not updated: %+v", second)
	}
	if second.ID != first.ID {
		t.Fatalf("upsert created a new row instead of updating")
	}
}

func TestUpsertProgressCompletedIsMonotonic(t *testing.T) {
	s := newTestStore(t)
	userID, exID := makeUserAndExercise(t, s)

	if _, err := s.UpsertProgress(userID, exID, "ok", true); err != nil {
		t.Fatalf("set completed: %v", err)
	}
	got, err := s.UpsertProgress(userID, exID, "broke it", false)
	if err != nil {
		t.Fatalf("re-upsert: %v", err)
	}
	if !got.Completed {
		t.Fatal("completed should stay true once an exercise is solved")
	}
}

func TestListProgressReturnsUserRows(t *testing.T) {
	s := newTestStore(t)
	userID, exID := makeUserAndExercise(t, s)
	_, _ = s.UpsertProgress(userID, exID, "x", false)

	list, err := s.ListProgress(userID)
	if err != nil {
		t.Fatalf("ListProgress: %v", err)
	}
	if len(list) != 1 || list[0].ExerciseID != exID {
		t.Fatalf("unexpected list: %+v", list)
	}
}
