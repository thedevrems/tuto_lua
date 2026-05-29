// Package payment integrates Stripe Checkout: it creates payment sessions and
// processes the webhook that unlocks a course once payment succeeds.
package payment

import (
	"encoding/json"
	"errors"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/webhook"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

// Service-level errors surfaced to handlers.
var (
	ErrNotConfigured   = errors.New("paiement non configuré")
	ErrFreeCourse      = errors.New("ce cours est gratuit")
	ErrAlreadyEnrolled = errors.New("vous avez déjà accès à ce cours")
)

// Store is the data surface the payment service relies on.
type Store interface {
	GetCourseByID(id string) (models.Course, error)
	IsEnrolled(userID, courseID string) (bool, error)
	CreatePayment(userID, courseID, sessionID string, amountCents int, currency string) (string, error)
	MarkPaymentPaid(sessionID string) (models.Payment, error)
	CreateEnrollment(userID, courseID, source string, grantedBy *string) error
}

// Service holds Stripe configuration and the data store.
type Service struct {
	store         Store
	frontendURL   string
	webhookSecret string
	enabled       bool
}

// NewService configures Stripe. When secretKey is empty the service is disabled
// and CreateCheckout returns ErrNotConfigured.
func NewService(store Store, secretKey, webhookSecret, frontendURL string) *Service {
	if secretKey != "" {
		stripe.Key = secretKey
	}
	return &Service{store: store, frontendURL: frontendURL, webhookSecret: webhookSecret, enabled: secretKey != ""}
}

// CreateCheckout starts a Stripe Checkout session for a user buying a course and
// returns the URL to redirect them to.
func (s *Service) CreateCheckout(userID, courseID string) (string, error) {
	if !s.enabled {
		return "", ErrNotConfigured
	}
	course, err := s.purchasableCourse(userID, courseID)
	if err != nil {
		return "", err
	}
	sess, err := session.New(s.checkoutParams(userID, course))
	if err != nil {
		return "", err
	}
	if _, err := s.store.CreatePayment(userID, course.ID, sess.ID, course.PriceCents, course.Currency); err != nil {
		return "", err
	}
	return sess.URL, nil
}

// purchasableCourse validates that the course exists, costs money and is not
// already owned by the user.
func (s *Service) purchasableCourse(userID, courseID string) (models.Course, error) {
	course, err := s.store.GetCourseByID(courseID)
	if err != nil {
		return models.Course{}, err
	}
	if course.PriceCents == 0 {
		return models.Course{}, ErrFreeCourse
	}
	enrolled, err := s.store.IsEnrolled(userID, courseID)
	if err != nil {
		return models.Course{}, err
	}
	if enrolled {
		return models.Course{}, ErrAlreadyEnrolled
	}
	return course, nil
}

// checkoutParams builds the Stripe session parameters for a one-off purchase.
func (s *Service) checkoutParams(userID string, course models.Course) *stripe.CheckoutSessionParams {
	return &stripe.CheckoutSessionParams{
		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:        stripe.String(s.frontendURL + "/learn?purchase=success"),
		CancelURL:         stripe.String(s.frontendURL + "/pricing?purchase=cancel"),
		ClientReferenceID: stripe.String(userID),
		Metadata:          map[string]string{"userId": userID, "courseId": course.ID},
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			Quantity: stripe.Int64(1),
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency:   stripe.String(course.Currency),
				UnitAmount: stripe.Int64(int64(course.PriceCents)),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(course.Title),
				},
			},
		}},
	}
}

// HandleWebhook verifies the Stripe signature and, on a completed checkout,
// marks the payment paid and grants course access.
func (s *Service) HandleWebhook(payload []byte, signature string) error {
	event, err := webhook.ConstructEvent(payload, signature, s.webhookSecret)
	if err != nil {
		return err
	}
	if event.Type != "checkout.session.completed" {
		return nil // ignore other event types
	}
	var sess stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
		return err
	}
	return s.fulfill(sess.ID)
}

// fulfill marks the payment paid and enrolls the buyer in the course.
func (s *Service) fulfill(sessionID string) error {
	paid, err := s.store.MarkPaymentPaid(sessionID)
	if err != nil {
		return err
	}
	return s.store.CreateEnrollment(paid.UserID, paid.CourseID, "purchase", nil)
}
