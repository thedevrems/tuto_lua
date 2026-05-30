// Package ticket orchestrates support reports and quote (devis) conversations:
// access control, message posting, closing, members and notifications.
package ticket

import (
	"errors"
	"strings"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

// Service-level errors surfaced to handlers.
var (
	ErrForbidden = errors.New("accès refusé à cette conversation")
	ErrEmpty     = errors.New("le contenu ne peut pas être vide")
)

// Store is the data surface the ticket service relies on.
type Store interface {
	CreateTicket(models.Ticket) (string, error)
	GetTicket(id string) (models.Ticket, error)
	ListUserTickets(userID string) ([]models.Ticket, error)
	ListAllTickets(ticketType string) ([]models.Ticket, error)
	CloseTicket(id string) error
	AddTicketMember(ticketID, userID string) error
	IsTicketMember(ticketID, userID string) (bool, error)
	ListTicketMembers(ticketID string) ([]models.TicketMember, error)
	ListTicketMemberIDs(ticketID string) ([]string, error)
	AddTicketMessage(ticketID, authorID, body string) (models.TicketMessage, error)
	ListTicketMessages(ticketID string) ([]models.TicketMessage, error)
	ListAdminIDs() ([]string, error)
	GetUserByID(id string) (models.User, error)
}

// Notifier fans an event out to a user's in-app inbox (and e-mail).
type Notifier interface {
	Notify(userID, title, body, link string)
}

// Service implements the ticket use-cases.
type Service struct {
	store    Store
	notifier Notifier
	link     string // base URL used in notification links
}

// NewService wires the store and notifier.
func NewService(store Store, notifier Notifier, frontendURL string) *Service {
	return &Service{store: store, notifier: notifier, link: frontendURL}
}

// CreateReport opens a support report with a first message and notifies admins.
func (s *Service) CreateReport(creator models.User, subject, category, body string) (models.Ticket, error) {
	subject, body = strings.TrimSpace(subject), strings.TrimSpace(body)
	if subject == "" || body == "" {
		return models.Ticket{}, ErrEmpty
	}
	id, err := s.store.CreateTicket(models.Ticket{
		Type: models.TicketReport, Subject: subject, Category: category,
		Status: models.TicketOpen, Currency: "eur", CreatorID: creator.ID,
	})
	if err != nil {
		return models.Ticket{}, err
	}
	if err := s.seedConversation(id, creator.ID, body); err != nil {
		return models.Ticket{}, err
	}
	s.notifyAdmins("Nouveau report — "+subject, creator.Username+" a ouvert un report.", id)
	return s.detail(id)
}

// seedConversation adds the creator as member and posts the opening message.
func (s *Service) seedConversation(ticketID, creatorID, body string) error {
	if err := s.store.AddTicketMember(ticketID, creatorID); err != nil {
		return err
	}
	_, err := s.store.AddTicketMessage(ticketID, creatorID, body)
	return err
}

// Detail returns a full ticket (messages + members) if the actor may see it.
func (s *Service) Detail(actor models.User, ticketID string) (models.Ticket, error) {
	t, err := s.store.GetTicket(ticketID)
	if err != nil {
		return models.Ticket{}, err
	}
	if !s.canAccess(actor, t) {
		return models.Ticket{}, ErrForbidden
	}
	return s.hydrate(t)
}

// PostMessage adds a reply (from a participant or an admin) and notifies others.
func (s *Service) PostMessage(actor models.User, ticketID, body string) (models.TicketMessage, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return models.TicketMessage{}, ErrEmpty
	}
	t, err := s.store.GetTicket(ticketID)
	if err != nil {
		return models.TicketMessage{}, err
	}
	if !s.canAccess(actor, t) {
		return models.TicketMessage{}, ErrForbidden
	}
	msg, err := s.store.AddTicketMessage(ticketID, actor.ID, body)
	if err != nil {
		return models.TicketMessage{}, err
	}
	s.notifyMembers(ticketID, actor.ID, "Nouveau message — "+t.Subject, actor.Username+" a répondu.")
	return msg, nil
}

// Close marks a ticket closed (admins only) and notifies its members.
func (s *Service) Close(actor models.User, ticketID string) error {
	t, err := s.requireAdminTicket(actor, ticketID)
	if err != nil {
		return err
	}
	if err := s.store.CloseTicket(ticketID); err != nil {
		return err
	}
	s.notifyMembers(ticketID, actor.ID, "Conversation fermée — "+t.Subject, "La conversation a été fermée.")
	return nil
}

// AddMember adds another participant to a ticket (admins only) and notifies them.
func (s *Service) AddMember(actor models.User, ticketID, userID string) error {
	t, err := s.requireAdminTicket(actor, ticketID)
	if err != nil {
		return err
	}
	if _, err := s.store.GetUserByID(userID); err != nil {
		return err
	}
	if err := s.store.AddTicketMember(ticketID, userID); err != nil {
		return err
	}
	s.notifier.Notify(userID, "Ajouté à une conversation — "+t.Subject, "Vous avez été ajouté à une conversation.", s.ticketLink(ticketID))
	return nil
}

// ListMine returns the tickets the user participates in.
func (s *Service) ListMine(userID string) ([]models.Ticket, error) {
	return s.store.ListUserTickets(userID)
}

// ListAll returns every ticket of a type (admins only; "" = all).
func (s *Service) ListAll(ticketType string) ([]models.Ticket, error) {
	return s.store.ListAllTickets(ticketType)
}

// requireAdminTicket checks the actor is admin and loads the ticket.
func (s *Service) requireAdminTicket(actor models.User, ticketID string) (models.Ticket, error) {
	if !actor.IsAdmin() {
		return models.Ticket{}, ErrForbidden
	}
	return s.store.GetTicket(ticketID)
}

// canAccess reports whether the actor may view/post in a ticket.
func (s *Service) canAccess(actor models.User, t models.Ticket) bool {
	if actor.IsAdmin() || actor.ID == t.CreatorID {
		return true
	}
	member, _ := s.store.IsTicketMember(t.ID, actor.ID)
	return member
}

// hydrate fills a ticket with its messages and members.
func (s *Service) hydrate(t models.Ticket) (models.Ticket, error) {
	var err error
	if t.Messages, err = s.store.ListTicketMessages(t.ID); err != nil {
		return models.Ticket{}, err
	}
	t.Members, err = s.store.ListTicketMembers(t.ID)
	return t, err
}

func (s *Service) detail(id string) (models.Ticket, error) {
	t, err := s.store.GetTicket(id)
	if err != nil {
		return models.Ticket{}, err
	}
	return s.hydrate(t)
}

// notifyMembers notifies every member except the actor.
func (s *Service) notifyMembers(ticketID, exceptID, title, body string) {
	ids, _ := s.store.ListTicketMemberIDs(ticketID)
	for _, id := range ids {
		if id != exceptID {
			s.notifier.Notify(id, title, body, s.ticketLink(ticketID))
		}
	}
}

// notifyAdmins notifies every administrator.
func (s *Service) notifyAdmins(title, body, ticketID string) {
	ids, _ := s.store.ListAdminIDs()
	for _, id := range ids {
		s.notifier.Notify(id, title, body, s.ticketLink(ticketID))
	}
}

func (s *Service) ticketLink(ticketID string) string {
	return s.link + "/support?t=" + ticketID
}
