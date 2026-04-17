package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type EmailSender interface {
	Send(ctx context.Context, to []string, subject, htmlBody string) error
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

// with dialer we can log where it fails and the req doesnt get stuck
func (s *SMTPSender) Send(ctx context.Context, to []string, subject string, htmlBody string) error {
	// skip in tests
	if os.Getenv("SKIP_EMAIL_SENDING") == "true" {
		log.Printf("Test mode: skipping email send to %v (subject: %s)", to, subject)
		return nil
	}

	if len(to) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	address := fmt.Sprintf("%s:%s", s.host, s.port)

	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}

	c, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("smtp client failed: %w", err)
	}
	defer c.Close()

	tlsConfig := &tls.Config{
		ServerName: s.host,
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}
	}

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}

	if err = c.Mail(s.username); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		htmlBody

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
