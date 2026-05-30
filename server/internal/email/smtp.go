// Package email sends transactional e-mails over SMTP. When no SMTP host is
// configured the mailer is a no-op, so the app runs fine without e-mail.
package email

import (
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

// SMTP sends mail through a configured SMTP server.
type SMTP struct {
	addr    string // host:port
	user    string
	pass    string
	from    string // may be "Name <addr>"
	enabled bool
}

// NewSMTP builds a mailer; it is disabled (no-op) when host is empty.
func NewSMTP(host, port, user, pass, from string) *SMTP {
	return &SMTP{addr: net.JoinHostPort(host, port), user: user, pass: pass, from: from, enabled: host != ""}
}

// Send delivers a plain-text e-mail. It returns nil immediately when disabled.
func (m *SMTP) Send(to, subject, body string) error {
	if !m.enabled {
		return nil
	}
	host, _, err := net.SplitHostPort(m.addr)
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth("", m.user, m.pass, host)
	return smtp.SendMail(m.addr, auth, addressOf(m.from), []string{to}, m.message(to, subject, body))
}

// message builds an RFC 822 plain-text message with CRLF line endings.
func (m *SMTP) message(to, subject, body string) []byte {
	headers := []string{
		"From: " + m.from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
	}
	return []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + body)
}

// addressOf extracts the bare e-mail from a possibly "Name <addr>" string.
func addressOf(from string) string {
	if parsed, err := mail.ParseAddress(from); err == nil {
		return parsed.Address
	}
	return from
}
