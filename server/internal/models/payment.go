package models

import "time"

// Payment statuses track a Stripe Checkout session through its lifecycle.
const (
	PaymentPending = "pending"
	PaymentPaid    = "paid"
	PaymentFailed  = "failed"
)

// Payment links a Stripe Checkout session to the course it unlocks.
type Payment struct {
	ID              string    `json:"id"`
	UserID          string    `json:"userId"`
	CourseID        string    `json:"courseId"`
	StripeSessionID string    `json:"stripeSessionId"`
	AmountCents     int       `json:"amountCents"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
}
