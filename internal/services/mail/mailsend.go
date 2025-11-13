package mail

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type MailSender struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewMailSender(host, port, username, password string) *MailSender {
	ms := &MailSender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}

	if ms.Host == "" {
		ms.Host = "smtp.gmail.com"
	}
	if ms.Port == "" {
		ms.Port = "587"
	}
	return ms
}

func SendMail(ms *MailSender, receiever []string, subject string, body string) error {
	if os.Getenv("SKIP_EMAIL_SENDING") == "true" {
		log.Printf("Test mode: skipping email send to %v (subject: %s)", receiever, subject)
		return nil
	}

	if len(receiever) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	address := fmt.Sprintf("%s:%s", ms.Host, ms.Port)
	auth := smtp.PlainAuth("", ms.Username, ms.Password, ms.Host)
	msg := []byte("To: " + strings.Join(receiever, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	if err := smtp.SendMail(address, auth, ms.Username, receiever, msg); err != nil {
		return err
	}

	return nil
}
