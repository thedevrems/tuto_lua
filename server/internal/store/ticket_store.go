package store

import (
	"database/sql"
	"errors"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

// ticketSelect joins the creator's username so listings can show who opened it.
const ticketSelect = `SELECT t.id, t.type, t.subject, t.category, t.status, t.details,
	t.amount_cents, t.currency, t.creator_id, t.created_at, t.updated_at, t.closed_at, u.username
	FROM tickets t JOIN users u ON u.id = t.creator_id`

// CreateTicket inserts a ticket (report or devis) and returns its id.
func (s *Store) CreateTicket(t models.Ticket) (string, error) {
	id := newID()
	_, err := s.db.Exec(
		`INSERT INTO tickets (id, type, subject, category, status, details, amount_cents, currency, creator_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, t.Type, t.Subject, t.Category, t.Status, t.Details, t.AmountCents, t.Currency, t.CreatorID)
	return id, err
}

// GetTicket loads a single ticket header (no messages/members).
func (s *Store) GetTicket(id string) (models.Ticket, error) {
	return scanTicket(s.db.QueryRow(ticketSelect+` WHERE t.id = ?`, id))
}

// ListUserTickets returns tickets the user created or is a member of, newest first.
func (s *Store) ListUserTickets(userID string) ([]models.Ticket, error) {
	return s.queryTickets(
		ticketSelect+` WHERE t.creator_id = ?
		 OR t.id IN (SELECT ticket_id FROM ticket_members WHERE user_id = ?)
		 ORDER BY t.updated_at DESC`, userID, userID)
}

// ListAllTickets returns every ticket, optionally filtered by type ("" = all).
func (s *Store) ListAllTickets(ticketType string) ([]models.Ticket, error) {
	if ticketType == "" {
		return s.queryTickets(ticketSelect + ` ORDER BY t.updated_at DESC`)
	}
	return s.queryTickets(ticketSelect+` WHERE t.type = ? ORDER BY t.updated_at DESC`, ticketType)
}

// CloseTicket marks a ticket closed with a timestamp.
func (s *Store) CloseTicket(id string) error {
	return s.execAffecting(
		`UPDATE tickets SET status = 'closed', closed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
}

// SetTicketStatus updates a ticket's status (used for devis accepted/refused).
func (s *Store) SetTicketStatus(id, status string) error {
	return s.execAffecting(`UPDATE tickets SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
}

// SetTicketAmount stores a devis amount (in cents) on the ticket.
func (s *Store) SetTicketAmount(id string, amountCents int) error {
	return s.execAffecting(`UPDATE tickets SET amount_cents = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, amountCents, id)
}

// queryTickets runs a ticketSelect query and scans the rows.
func (s *Store) queryTickets(query string, args ...any) ([]models.Ticket, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		t, err := scanTicket(rows)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, rows.Err()
}

// scanTicket reads one ticket row (with creator username and nullable closed_at).
func scanTicket(row rowScanner) (models.Ticket, error) {
	var t models.Ticket
	var closed sql.NullTime
	err := row.Scan(&t.ID, &t.Type, &t.Subject, &t.Category, &t.Status, &t.Details,
		&t.AmountCents, &t.Currency, &t.CreatorID, &t.CreatedAt, &t.UpdatedAt, &closed, &t.CreatorName)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Ticket{}, ErrNotFound
	}
	if err != nil {
		return models.Ticket{}, err
	}
	if closed.Valid {
		t.ClosedAt = &closed.Time
	}
	return t, nil
}
