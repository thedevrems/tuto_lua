// Package notify creates in-app notifications and mirrors them by e-mail.
package notify

import "github.com/thedevrems/tuto_lua/server/internal/models"

// Store persists notifications and resolves a user's e-mail.
type Store interface {
	CreateNotification(userID, title, body, link string) error
	GetUserByID(id string) (models.User, error)
}

// Mailer sends an e-mail (a no-op implementation is fine when unconfigured).
type Mailer interface {
	Send(to, subject, body string) error
}

// Notifier fans a single event out to the in-app inbox and to e-mail.
type Notifier struct {
	store  Store
	mailer Mailer
}

// New builds a Notifier from the store and mailer.
func New(store Store, mailer Mailer) *Notifier {
	return &Notifier{store: store, mailer: mailer}
}

// Notify records an in-app notification and e-mails the user. It is best-effort:
// any failure is swallowed so notifications never break the triggering action.
func (n *Notifier) Notify(userID, title, body, link string) {
	_ = n.store.CreateNotification(userID, title, body, link)
	if user, err := n.store.GetUserByID(userID); err == nil {
		_ = n.mailer.Send(user.Email, title, body)
	}
}
