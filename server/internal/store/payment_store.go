package store

import (
	"database/sql"
	"errors"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

const paymentColumns = `id, user_id, course_id, stripe_session_id, amount_cents, currency, status, created_at`

// CreatePayment records a pending payment tied to a Stripe Checkout session.
func (s *Store) CreatePayment(userID, courseID, sessionID string, amountCents int, currency string) (string, error) {
	id := newID()
	_, err := s.db.Exec(
		`INSERT INTO payments (id, user_id, course_id, stripe_session_id, amount_cents, currency, status)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, userID, courseID, sessionID, amountCents, currency, models.PaymentPending)
	return id, err
}

// MarkPaymentPaid flips a payment to "paid" and returns it (to read user/course).
func (s *Store) MarkPaymentPaid(sessionID string) (models.Payment, error) {
	if _, err := s.db.Exec(
		`UPDATE payments SET status = ? WHERE stripe_session_id = ?`, models.PaymentPaid, sessionID); err != nil {
		return models.Payment{}, err
	}
	return s.getPaymentBySession(sessionID)
}

// getPaymentBySession loads a payment by its Stripe session id.
func (s *Store) getPaymentBySession(sessionID string) (models.Payment, error) {
	row := s.db.QueryRow(`SELECT `+paymentColumns+` FROM payments WHERE stripe_session_id = ?`, sessionID)
	var p models.Payment
	err := row.Scan(&p.ID, &p.UserID, &p.CourseID, &p.StripeSessionID, &p.AmountCents, &p.Currency, &p.Status, &p.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Payment{}, ErrNotFound
	}
	if err != nil {
		return models.Payment{}, err
	}
	return p, nil
}
