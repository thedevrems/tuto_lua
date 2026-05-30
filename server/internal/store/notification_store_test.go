package store

import (
	"testing"

	"github.com/thedevrems/tuto_lua/server/internal/models"
)

func TestNotificationLifecycle(t *testing.T) {
	s := newTestStore(t)
	user, _ := s.CreateUser("u", "u@example.com", "h", models.RoleUser)

	if err := s.CreateNotification(user.ID, "Bienvenue", "Bonjour", "/profile"); err != nil {
		t.Fatalf("CreateNotification: %v", err)
	}
	if err := s.CreateNotification(user.ID, "Report fermé", "Votre report a été fermé", ""); err != nil {
		t.Fatalf("CreateNotification 2: %v", err)
	}

	list, err := s.ListNotifications(user.ID)
	if err != nil || len(list) != 2 {
		t.Fatalf("ListNotifications = (%d, %v), want 2", len(list), err)
	}
	if unread, _ := s.CountUnread(user.ID); unread != 2 {
		t.Fatalf("unread = %d, want 2", unread)
	}

	if err := s.MarkNotificationRead(user.ID, list[0].ID); err != nil {
		t.Fatalf("MarkNotificationRead: %v", err)
	}
	if unread, _ := s.CountUnread(user.ID); unread != 1 {
		t.Fatalf("unread after one read = %d, want 1", unread)
	}

	if err := s.MarkAllNotificationsRead(user.ID); err != nil {
		t.Fatalf("MarkAllNotificationsRead: %v", err)
	}
	if unread, _ := s.CountUnread(user.ID); unread != 0 {
		t.Fatalf("unread after all read = %d, want 0", unread)
	}
}
