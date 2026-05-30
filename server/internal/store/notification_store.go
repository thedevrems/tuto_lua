package store

import "github.com/thedevrems/tuto_lua/server/internal/models"

const notificationColumns = `id, user_id, title, body, link, read, created_at`

// CreateNotification stores an in-app notification for a user.
func (s *Store) CreateNotification(userID, title, body, link string) error {
	_, err := s.db.Exec(
		`INSERT INTO notifications (id, user_id, title, body, link) VALUES (?, ?, ?, ?, ?)`,
		newID(), userID, title, body, link)
	return err
}

// ListNotifications returns a user's notifications, newest first.
func (s *Store) ListNotifications(userID string) ([]models.Notification, error) {
	rows, err := s.db.Query(
		`SELECT `+notificationColumns+` FROM notifications WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Notification
	for rows.Next() {
		var n models.Notification
		var read int
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.Link, &read, &n.CreatedAt); err != nil {
			return nil, err
		}
		n.Read = read != 0
		list = append(list, n)
	}
	return list, rows.Err()
}

// CountUnread returns how many unread notifications a user has.
func (s *Store) CountUnread(userID string) (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = ? AND read = 0`, userID).Scan(&n)
	return n, err
}

// MarkNotificationRead marks one notification read (scoped to its owner).
func (s *Store) MarkNotificationRead(userID, id string) error {
	return s.execAffecting(`UPDATE notifications SET read = 1 WHERE id = ? AND user_id = ?`, id, userID)
}

// MarkAllNotificationsRead marks every notification of a user as read.
func (s *Store) MarkAllNotificationsRead(userID string) error {
	_, err := s.db.Exec(`UPDATE notifications SET read = 1 WHERE user_id = ?`, userID)
	return err
}
