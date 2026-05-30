package store

import (
	"testing"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

func TestTicketCreateMessageAndClose(t *testing.T) {
	s := newTestStore(t)
	u, _ := s.CreateUser("alice", "a@example.com", "h", models.RoleUser)

	id, err := s.CreateTicket(models.Ticket{
		Type: models.TicketReport, Subject: "Bug d'affichage", Category: "bug",
		Status: models.TicketOpen, Currency: "eur", CreatorID: u.ID,
	})
	if err != nil {
		t.Fatalf("CreateTicket: %v", err)
	}
	if err := s.AddTicketMember(id, u.ID); err != nil {
		t.Fatalf("AddTicketMember: %v", err)
	}
	msg, err := s.AddTicketMessage(id, u.ID, "ça plante")
	if err != nil || msg.AuthorName != "alice" || msg.Body != "ça plante" {
		t.Fatalf("AddTicketMessage = (%+v, %v)", msg, err)
	}

	got, err := s.GetTicket(id)
	if err != nil || got.CreatorName != "alice" || got.Subject != "Bug d'affichage" {
		t.Fatalf("GetTicket = (%+v, %v)", got, err)
	}
	if msgs, _ := s.ListTicketMessages(id); len(msgs) != 1 {
		t.Fatalf("messages = %d, want 1", len(msgs))
	}
	if mine, _ := s.ListUserTickets(u.ID); len(mine) != 1 {
		t.Fatalf("ListUserTickets = %d, want 1", len(mine))
	}
	if reports, _ := s.ListAllTickets(models.TicketReport); len(reports) != 1 {
		t.Fatalf("ListAllTickets(report) = %d, want 1", len(reports))
	}
	if devis, _ := s.ListAllTickets(models.TicketDevis); len(devis) != 0 {
		t.Fatalf("ListAllTickets(devis) = %d, want 0", len(devis))
	}

	if err := s.CloseTicket(id); err != nil {
		t.Fatalf("CloseTicket: %v", err)
	}
	closed, _ := s.GetTicket(id)
	if closed.Status != models.TicketClosed || closed.ClosedAt == nil {
		t.Fatalf("ticket not closed: %+v", closed)
	}
}

func TestAddTicketMemberIsIdempotent(t *testing.T) {
	s := newTestStore(t)
	u, _ := s.CreateUser("alice", "a@example.com", "h", models.RoleUser)
	id, _ := s.CreateTicket(models.Ticket{Type: models.TicketReport, Subject: "S", Status: models.TicketOpen, Currency: "eur", CreatorID: u.ID})

	_ = s.AddTicketMember(id, u.ID)
	if err := s.AddTicketMember(id, u.ID); err != nil {
		t.Fatalf("second AddTicketMember should be a no-op: %v", err)
	}
	members, _ := s.ListTicketMembers(id)
	if len(members) != 1 {
		t.Fatalf("members = %d, want 1", len(members))
	}
	if member, _ := s.IsTicketMember(id, u.ID); !member {
		t.Fatal("expected user to be a member")
	}
}
