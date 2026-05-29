package store

import (
	"testing"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

func TestEnrollmentGrantAndCheck(t *testing.T) {
	s := newTestStore(t)
	user, _ := s.CreateUser("buyer", "b@example.com", "h", models.RoleUser)
	courseID, _ := s.CreateCourse(models.Course{Slug: "paid", Title: "Paid", Currency: "eur", PriceCents: 4900})

	enrolled, err := s.IsEnrolled(user.ID, courseID)
	if err != nil || enrolled {
		t.Fatalf("expected not enrolled, got (%v, %v)", enrolled, err)
	}

	if err := s.CreateEnrollment(user.ID, courseID, "purchase", nil); err != nil {
		t.Fatalf("CreateEnrollment: %v", err)
	}
	enrolled, _ = s.IsEnrolled(user.ID, courseID)
	if !enrolled {
		t.Fatal("expected enrolled after grant")
	}
}

func TestCreateEnrollmentIsIdempotent(t *testing.T) {
	s := newTestStore(t)
	user, _ := s.CreateUser("buyer", "b@example.com", "h", models.RoleUser)
	courseID, _ := s.CreateCourse(models.Course{Slug: "paid", Title: "Paid", Currency: "eur", PriceCents: 4900})

	if err := s.CreateEnrollment(user.ID, courseID, "purchase", nil); err != nil {
		t.Fatalf("first grant: %v", err)
	}
	if err := s.CreateEnrollment(user.ID, courseID, "admin", nil); err != nil {
		t.Fatalf("second grant should be a no-op, got: %v", err)
	}
	ids, _ := s.ListEnrollmentCourseIDs(user.ID)
	if len(ids) != 1 || ids[0] != courseID {
		t.Fatalf("expected exactly one enrollment, got %+v", ids)
	}
}
