package mail

import (
	"fmt"
	"net/smtp"
)

type MailSender struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (ms *MailSender) SendMail(reciever []string, subject string, body string) error {
	if len(reciever) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	address := fmt.Sprintf("%s:%s", ms.Host, ms.Port)
	auth := smtp.PlainAuth("", ms.Username, ms.Password, ms.Host)
	msg := []byte("To: " + reciever[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	if err := smtp.SendMail(address, auth, ms.Username, reciever, msg); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
