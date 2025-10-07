package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

type MailSender struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (ms *MailSender) SendMail(receiever []string, subject string, body string) error {
	if len(receiever) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	address := fmt.Sprintf("%s:%s", ms.Host, ms.Port)
	auth := smtp.PlainAuth("", ms.Username, ms.Password, ms.Host)
	msg := []byte("To: " + strings.Join(receiever, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	if err := smtp.SendMail(address, auth, ms.Username, receiever, msg); err != nil {
		return err
	}

	return nil
}
