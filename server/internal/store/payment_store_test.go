package store

import (
	"errors"
	"testing"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

func TestCreatePaymentThenMarkPaid(t *testing.T) {
	s := newTestStore(t)
	user, _ := s.CreateUser("buyer", "b@example.com", "h", models.RoleUser)
	courseID, _ := s.CreateCourse(models.Course{Slug: "p", Title: "Paid", Currency: "eur", PriceCents: 4900})

	id, err := s.CreatePayment(user.ID, courseID, "sess_123", 4900, "eur")
	if err != nil || id == "" {
		t.Fatalf("CreatePayment = (%q, %v)", id, err)
	}

	paid, err := s.MarkPaymentPaid("sess_123")
	if err != nil {
		t.Fatalf("MarkPaymentPaid: %v", err)
	}
	if paid.Status != models.PaymentPaid {
		t.Fatalf("status = %q, want paid", paid.Status)
	}
	if paid.UserID != user.ID || paid.CourseID != courseID {
		t.Fatalf("payment links wrong: %+v", paid)
	}
}

func TestMarkPaymentPaidUnknownSession(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.MarkPaymentPaid("does-not-exist"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
