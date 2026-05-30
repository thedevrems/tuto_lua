package store

import (
	"database/sql"
	"errors"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

// userColumns is the canonical SELECT list, reused by every user query.
const userColumns = `id, username, email, password_hash, role, created_at`

// CreateUser inserts a new account and returns it with its generated id.
// It maps a UNIQUE violation (username/email taken) to ErrConflict.
func (s *Store) CreateUser(username, email, passwordHash string, role models.Role) (models.User, error) {
	u := models.User{ID: newID(), Username: username, Email: email, PasswordHash: passwordHash, Role: role}
	_, err := s.db.Exec(
		`INSERT INTO users (id, username, email, password_hash, role) VALUES (?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.Email, u.PasswordHash, string(u.Role),
	)
	if isUniqueViolation(err) {
		return models.User{}, ErrConflict
	}
	if err != nil {
		return models.User{}, err
	}
	return s.GetUserByID(u.ID)
}

// GetUserByID loads a single account by primary key.
func (s *Store) GetUserByID(id string) (models.User, error) {
	return scanUser(s.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE id = ?`, id))
}

// GetUserByEmail loads an account by its (case-insensitive) email.
func (s *Store) GetUserByEmail(email string) (models.User, error) {
	return scanUser(s.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE email = ? COLLATE NOCASE`, email))
}

// GetUserByUsername loads an account by its (case-insensitive) handle.
func (s *Store) GetUserByUsername(username string) (models.User, error) {
	return scanUser(s.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE username = ? COLLATE NOCASE`, username))
}

// ListAdminIDs returns the ids of every administrator (to notify them).
func (s *Store) ListAdminIDs() ([]string, error) {
	rows, err := s.db.Query(`SELECT id FROM users WHERE role = 'admin'`)
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

// UpdatePassword replaces a user's password hash.
func (s *Store) UpdatePassword(userID, passwordHash string) error {
	return s.execAffecting(`UPDATE users SET password_hash = ? WHERE id = ?`, passwordHash, userID)
}

// CountUsers returns the total number of accounts (used to bootstrap admin).
func (s *Store) CountUsers() (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&n)
	return n, err
}

// ListUsers returns every account, oldest first (admin user management).
func (s *Store) ListUsers() ([]models.User, error) {
	rows, err := s.db.Query(`SELECT ` + userColumns + ` FROM users ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// scanUser reads one user row, translating sql.ErrNoRows into ErrNotFound.
func scanUser(row rowScanner) (models.User, error) {
	var u models.User
	var role string
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &role, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	u.Role = models.Role(role)
	return u, nil
}
