package ticket_test

import (
	"errors"
	"testing"

	"github.com/thedevrems/tuto_lua/server/internal/database"
	"github.com/thedevrems/tuto_lua/server/internal/models"
	"github.com/thedevrems/tuto_lua/server/internal/store"
	"github.com/thedevrems/tuto_lua/server/internal/ticket"
)

// recNotifier records who was notified.
type recNotifier struct{ users []string }

func (r *recNotifier) Notify(userID, _, _, _ string) { r.users = append(r.users, userID) }
func (r *recNotifier) notified(id string) bool {
	for _, u := range r.users {
		if u == id {
			return true
		}
	}
	return false
}

func setup(t *testing.T) (*ticket.Service, *recNotifier, *store.Store) {
	t.Helper()
	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	st := store.New(db)
	n := &recNotifier{}
	return ticket.NewService(st, n, "http://front"), n, st
}

func TestCreateReportNotifiesAdmins(t *testing.T) {
	svc, notifier, st := setup(t)
	admin, _ := st.CreateUser("boss", "boss@e.com", "h", models.RoleAdmin)
	user, _ := st.CreateUser("alice", "a@e.com", "h", models.RoleUser)

	tk, err := svc.CreateReport(user, "Bug", "bug", "ça plante")
	if err != nil {
		t.Fatalf("CreateReport: %v", err)
	}
	if len(tk.Messages) != 1 || len(tk.Members) != 1 {
		t.Fatalf("expected 1 message + 1 member, got %+v", tk)
	}
	if !notifier.notified(admin.ID) {
		t.Fatal("admin should be notified of a new report")
	}
}

func TestCreateReportRejectsEmpty(t *testing.T) {
	svc, _, st := setup(t)
	user, _ := st.CreateUser("alice", "a@e.com", "h", models.RoleUser)
	if _, err := svc.CreateReport(user, "", "bug", ""); !errors.Is(err, ticket.ErrEmpty) {
		t.Fatalf("expected ErrEmpty, got %v", err)
	}
}

func TestAccessControl(t *testing.T) {
	svc, _, st := setup(t)
	admin, _ := st.CreateUser("boss", "boss@e.com", "h", models.RoleAdmin)
	owner, _ := st.CreateUser("alice", "a@e.com", "h", models.RoleUser)
	outsider, _ := st.CreateUser("bob", "b@e.com", "h", models.RoleUser)

	tk, _ := svc.CreateReport(owner, "Bug", "bug", "aide")

	if _, err := svc.Detail(owner, tk.ID); err != nil {
		t.Fatalf("owner should access: %v", err)
	}
	if _, err := svc.Detail(admin, tk.ID); err != nil {
		t.Fatalf("admin should access: %v", err)
	}
	if _, err := svc.Detail(outsider, tk.ID); !errors.Is(err, ticket.ErrForbidden) {
		t.Fatalf("outsider should be forbidden, got %v", err)
	}
}

func TestCloseRequiresAdminAndNotifies(t *testing.T) {
	svc, notifier, st := setup(t)
	_, _ = st.CreateUser("boss", "boss@e.com", "h", models.RoleAdmin)
	admin, _ := st.GetUserByUsername("boss")
	owner, _ := st.CreateUser("alice", "a@e.com", "h", models.RoleUser)

	tk, _ := svc.CreateReport(owner, "Bug", "bug", "aide")

	if err := svc.Close(owner, tk.ID); !errors.Is(err, ticket.ErrForbidden) {
		t.Fatalf("non-admin close should be forbidden, got %v", err)
	}
	notifier.users = nil
	if err := svc.Close(admin, tk.ID); err != nil {
		t.Fatalf("admin close: %v", err)
	}
	if !notifier.notified(owner.ID) {
		t.Fatal("owner should be notified when the ticket is closed")
	}
}

func TestAddMemberGrantsAccess(t *testing.T) {
	svc, notifier, st := setup(t)
	admin, _ := st.CreateUser("boss", "boss@e.com", "h", models.RoleAdmin)
	owner, _ := st.CreateUser("alice", "a@e.com", "h", models.RoleUser)
	other, _ := st.CreateUser("bob", "b@e.com", "h", models.RoleUser)

	tk, _ := svc.CreateReport(owner, "Bug", "bug", "aide")
	if err := svc.AddMember(admin, tk.ID, other.ID); err != nil {
		t.Fatalf("AddMember: %v", err)
	}
	if !notifier.notified(other.ID) {
		t.Fatal("added member should be notified")
	}
	if _, err := svc.Detail(other, tk.ID); err != nil {
		t.Fatalf("added member should now access: %v", err)
	}
}
