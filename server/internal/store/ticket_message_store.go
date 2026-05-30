package store

import "github.com/thedevrems/tuto_lua/server/internal/models"

// AddTicketMember adds a participant to a ticket, idempotently.
func (s *Store) AddTicketMember(ticketID, userID string) error {
	_, err := s.db.Exec(
		`INSERT INTO ticket_members (ticket_id, user_id) VALUES (?, ?) ON CONFLICT DO NOTHING`, ticketID, userID)
	return err
}

// IsTicketMember reports whether a user participates in a ticket.
func (s *Store) IsTicketMember(ticketID, userID string) (bool, error) {
	var one int
	err := s.db.QueryRow(
		`SELECT 1 FROM ticket_members WHERE ticket_id = ? AND user_id = ? LIMIT 1`, ticketID, userID).Scan(&one)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ListTicketMembers returns a ticket's participants with their usernames.
func (s *Store) ListTicketMembers(ticketID string) ([]models.TicketMember, error) {
	rows, err := s.db.Query(
		`SELECT m.user_id, u.username FROM ticket_members m
		 JOIN users u ON u.id = m.user_id WHERE m.ticket_id = ? ORDER BY m.added_at`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TicketMember
	for rows.Next() {
		var m models.TicketMember
		if err := rows.Scan(&m.UserID, &m.Username); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

// ListTicketMemberIDs returns just the user ids participating in a ticket
// (used to fan out notifications).
func (s *Store) ListTicketMemberIDs(ticketID string) ([]string, error) {
	rows, err := s.db.Query(`SELECT user_id FROM ticket_members WHERE ticket_id = ?`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// AddTicketMessage posts a message and bumps the ticket's updated_at.
func (s *Store) AddTicketMessage(ticketID, authorID, body string) (models.TicketMessage, error) {
	id := newID()
	if _, err := s.db.Exec(
		`INSERT INTO ticket_messages (id, ticket_id, author_id, body) VALUES (?, ?, ?, ?)`,
		id, ticketID, authorID, body); err != nil {
		return models.TicketMessage{}, err
	}
	if _, err := s.db.Exec(`UPDATE tickets SET updated_at = CURRENT_TIMESTAMP WHERE id = ?`, ticketID); err != nil {
		return models.TicketMessage{}, err
	}
	return s.getMessage(id)
}

// ListTicketMessages returns a ticket's messages oldest-first with author names.
func (s *Store) ListTicketMessages(ticketID string) ([]models.TicketMessage, error) {
	rows, err := s.db.Query(messageSelect+` WHERE m.ticket_id = ? ORDER BY m.created_at`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.TicketMessage
	for rows.Next() {
		m, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

const messageSelect = `SELECT m.id, m.ticket_id, m.author_id, u.username, m.body, m.created_at
	FROM ticket_messages m JOIN users u ON u.id = m.author_id`

// getMessage loads a single message by id (with author name).
func (s *Store) getMessage(id string) (models.TicketMessage, error) {
	return scanMessage(s.db.QueryRow(messageSelect+` WHERE m.id = ?`, id))
}

func scanMessage(row rowScanner) (models.TicketMessage, error) {
	var m models.TicketMessage
	if err := row.Scan(&m.ID, &m.TicketID, &m.AuthorID, &m.AuthorName, &m.Body, &m.CreatedAt); err != nil {
		return models.TicketMessage{}, err
	}
	return m, nil
}
