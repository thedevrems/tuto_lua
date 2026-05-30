package models

import "time"

// Ticket types and statuses. A ticket is a conversation: either a support
// "report" or a "devis" (quote request).
const (
	TicketReport = "report"
	TicketDevis  = "devis"

	TicketOpen     = "open"
	TicketClosed   = "closed"
	TicketAccepted = "accepted" // devis only
	TicketRefused  = "refused"  // devis only
)

// Ticket is a conversation thread shared by support reports and quotes.
type Ticket struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Subject     string     `json:"subject"`
	Category    string     `json:"category"`
	Status      string     `json:"status"`
	Details     string     `json:"details"` // JSON payload for devis form data
	AmountCents int        `json:"amountCents"`
	Currency    string     `json:"currency"`
	CreatorID   string     `json:"creatorId"`
	CreatorName string     `json:"creatorName,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ClosedAt    *time.Time `json:"closedAt,omitempty"`
	// Populated only when a full ticket is requested.
	Messages []TicketMessage `json:"messages,omitempty"`
	Members  []TicketMember  `json:"members,omitempty"`
}

// TicketMessage is a single post in a ticket conversation.
type TicketMessage struct {
	ID         string    `json:"id"`
	TicketID   string    `json:"ticketId"`
	AuthorID   string    `json:"authorId"`
	AuthorName string    `json:"authorName"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"createdAt"`
}

// TicketMember is a participant of a ticket (creator + added members).
type TicketMember struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
}
