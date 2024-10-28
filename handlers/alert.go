package handlers

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"goapi/dbconnect"
	"goapi/models"
	"log"
	"os"
	"time"
)

func SendAlert() {
	mailsend := models.MailSender{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     os.Getenv("EMAIL_PORT"),
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}
	var name, email string
	var net_amount, gross_amount float64
	var listPayments = make(map[string]string)

	db := dbconnect.DB
	if db == nil {
		log.Fatal("could not initialize the db")
	}
	tomorrow := time.Now().Add(24 * time.Hour)

	rows, err := db.Table("users").Select("users.email, payments.name, payments.net_amount, payments.gross_amount").Joins("left join payments on payments.user_id = users.id").Where("payments.date = ?", tomorrow).Rows()
	if err != nil {
		log.Fatal(err)
	}

	body_massage := "these are the payments that will be executed tomorrow. \r\n"
	for rows.Next() {
		rows.Scan(&name, &email, &net_amount, &gross_amount)
		details := fmt.Sprintf("payment: %s \r\n net amount: %.2f \r\n gross amount: %.2f \r\n", name, net_amount, gross_amount)
		listPayments[email] = details
	}

	for key, val := range listPayments {
		mailsend.SendMail([]string{key}, "pay alert", body_massage+val)

	}
}
