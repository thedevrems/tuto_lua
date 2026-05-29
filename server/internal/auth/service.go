package auth

import (
	"errors"
	"strings"

	"github.com/thedevrems/tuto_lua/server/internal/crypto"
	"github.com/thedevrems/tuto_lua/server/internal/models"
	"github.com/thedevrems/tuto_lua/server/internal/token"
	"github.com/thedevrems/tuto_lua/server/internal/validate"
)

// ErrInvalidCredentials is returned for any failed login. It is deliberately
// vague so attackers cannot tell whether the account exists.
var ErrInvalidCredentials = errors.New("identifiants invalides")

// UserStore is the subset of the data layer the auth service depends on.
type UserStore interface {
	CreateUser(username, email, passwordHash string, role models.Role) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	CountUsers() (int, error)
}

// Service turns credentials into users and signed tokens.
type Service struct {
	users  UserStore
	tokens *token.Manager
}

// Result is the payload returned to clients after auth succeeds.
type Result struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

// NewService wires the store and token manager together.
func NewService(users UserStore, tokens *token.Manager) *Service {
	return &Service{users: users, tokens: tokens}
}

// Register validates the input, creates the account and issues a token.
// The very first account created becomes an admin to bootstrap the platform.
func (s *Service) Register(username, email, password string) (Result, error) {
	username, email, err := validateRegistration(username, email, password)
	if err != nil {
		return Result{}, err
	}
	hash, err := crypto.HashPassword(password)
	if err != nil {
		return Result{}, err
	}
	user, err := s.users.CreateUser(username, email, hash, s.firstUserRole())
	if err != nil {
		return Result{}, err
	}
	return s.issue(user)
}

// Login authenticates by email or username and returns a fresh token.
func (s *Service) Login(identifier, password string) (Result, error) {
	user, err := s.lookup(identifier)
	if err != nil {
		return Result{}, ErrInvalidCredentials
	}
	if err := crypto.VerifyPassword(user.PasswordHash, password); err != nil {
		return Result{}, ErrInvalidCredentials
	}
	return s.issue(user)
}

// lookup finds a user by email when the identifier looks like one, else by name.
func (s *Service) lookup(identifier string) (models.User, error) {
	id := strings.TrimSpace(identifier)
	if strings.Contains(id, "@") {
		return s.users.GetUserByEmail(strings.ToLower(id))
	}
	return s.users.GetUserByUsername(id)
}

// firstUserRole returns admin for the bootstrap account, user otherwise.
func (s *Service) firstUserRole() models.Role {
	if n, err := s.users.CountUsers(); err == nil && n == 0 {
		return models.RoleAdmin
	}
	return models.RoleUser
}

// issue mints a token for an authenticated user.
func (s *Service) issue(user models.User) (Result, error) {
	tok, err := s.tokens.Issue(user.ID, string(user.Role))
	if err != nil {
		return Result{}, err
	}
	return Result{User: user, Token: tok}, nil
}

// validateRegistration normalizes and checks all sign-up fields up front.
func validateRegistration(username, email, password string) (string, string, error) {
	u, err := validate.Username(username)
	if err != nil {
		return "", "", err
	}
	e, err := validate.Email(email)
	if err != nil {
		return "", "", err
	}
	if err := validate.Password(password); err != nil {
		return "", "", err
	}
	return u, e, nil
}
