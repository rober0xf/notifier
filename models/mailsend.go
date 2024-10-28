package models

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

type MailSender struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (ms *MailSender) SendMail(reciever []string, subject string, body string) error {
	address := fmt.Sprintf("%s:%s", ms.Host, ms.Port)
	tls_config := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         ms.Host,
	}

	connection, err := tls.Dial("tcp", address, tls_config)
	if err != nil {
		log.Fatal(err)
		return err
	}

	client, err := smtp.NewClient(connection, ms.Host)
	if err != nil {
		log.Fatal(err)
		return err
	}

	auth := smtp.PlainAuth("", ms.Username, ms.Password, ms.Host)
	if err = client.Auth(auth); err != nil {
		log.Fatal(err)
		return err
	}

	if err = client.Mail(ms.Username); err != nil {
		log.Fatal(err)
		return err
	}

	for _, r := range reciever {
		if err = client.Rcpt(r); err != nil {
			log.Fatal(err)
			return err
		}
	}

	wc, err := client.Data()
	if err != nil {
		log.Fatal(err)
		return err
	}

	message := []byte("to: " + reciever[0] + "\r\n" +
		"subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	_, err = wc.Write(message)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = wc.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}

	client.Quit()
	/*
		err := smtp.SendMail(address, auth, ms.Username, reciever, message)
		if err != nil {
			return err
		}
	*/
	return nil
}
