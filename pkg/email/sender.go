package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type EmailSender interface {
	Send(to []string, subject, htmlBody string) error
}

type SMTPSender struct {
	host     string
	port     string
	username string
	password string
}

func NewSMTPSender(host, port, username, password string) *SMTPSender {
	if host == "" {
		host = "smtp.gmail.com"
	}
	if port == "" {
		port = "587"
	}

	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (s *SMTPSender) Send(to []string, subject string, htmlBody string) error {
	// skip in tests
	if os.Getenv("SKIP_EMAIL_SENDING") == "true" {
		log.Printf("Test mode: skipping email send to %v (subject: %s)", to, subject)
		return nil
	}

	if len(to) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	address := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	msg := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		htmlBody + "\r\n")

	if err := smtp.SendMail(address, auth, s.username, to, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
