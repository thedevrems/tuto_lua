package payment

import (
	"errors"
	"testing"

	"github.com/thedevrems/tuto_lua/server/internal/models"
	"github.com/thedevrems/tuto_lua/server/internal/store"
)

// fakeStore lets us exercise the pre-Stripe validation paths in isolation.
type fakeStore struct {
	course   models.Course
	notFound bool
	enrolled bool
}

func (f *fakeStore) GetCourseByID(string) (models.Course, error) {
	if f.notFound {
		return models.Course{}, store.ErrNotFound
	}
	return f.course, nil
}
func (f *fakeStore) IsEnrolled(string, string) (bool, error)             { return f.enrolled, nil }
func (f *fakeStore) CreatePayment(string, string, string, int, string) (string, error) { return "p1", nil }
func (f *fakeStore) MarkPaymentPaid(string) (models.Payment, error)      { return models.Payment{}, nil }
func (f *fakeStore) CreateEnrollment(string, string, string, *string) error { return nil }

func TestCreateCheckoutNotConfigured(t *testing.T) {
	svc := NewService(&fakeStore{}, "", "", "http://front")
	if _, err := svc.CreateCheckout("u1", "c1"); !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}
}

func TestCreateCheckoutRejectsFreeCourse(t *testing.T) {
	fs := &fakeStore{course: models.Course{ID: "c1", PriceCents: 0, Currency: "eur"}}
	svc := NewService(fs, "sk_test_x", "", "http://front")
	if _, err := svc.CreateCheckout("u1", "c1"); !errors.Is(err, ErrFreeCourse) {
		t.Fatalf("expected ErrFreeCourse, got %v", err)
	}
}

func TestCreateCheckoutRejectsAlreadyEnrolled(t *testing.T) {
	fs := &fakeStore{course: models.Course{ID: "c1", PriceCents: 4900, Currency: "eur"}, enrolled: true}
	svc := NewService(fs, "sk_test_x", "", "http://front")
	if _, err := svc.CreateCheckout("u1", "c1"); !errors.Is(err, ErrAlreadyEnrolled) {
		t.Fatalf("expected ErrAlreadyEnrolled, got %v", err)
	}
}

func TestCreateCheckoutCourseNotFound(t *testing.T) {
	svc := NewService(&fakeStore{notFound: true}, "sk_test_x", "", "http://front")
	if _, err := svc.CreateCheckout("u1", "missing"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected store.ErrNotFound, got %v", err)
	}
}
