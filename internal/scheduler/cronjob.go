package cronjob

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	_ "github.com/joho/godotenv/autoload"
	database "github.com/rober0xf/notifier/internal/db"
	"github.com/rober0xf/notifier/internal/models"
)

func InitCron() {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("error creating scheduler: %v", err)
		return
	}

	job, err := s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(
			gocron.NewAtTime(15, 27, 00),
		),
		),
		gocron.NewTask(SendDailyAlert),
	)
	if err != nil {
		log.Print("error creating new job: ", err)
	}

	log.Printf("ID: %v", job.ID())
	s.Start()
}

func SendDailyAlert() {
	mailsend := models.MailSender{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}

	var email, name string
	var net_amount, gross_amount float64
	var listPayments = make(map[string]string)

	db := database.DB
	if db == nil {
		log.Fatal("could not initialize database")
	}

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	var count int64

	rows, err := db.Table("users").
		Select("users.email, payments.name, payments.net_amount, payments.gross_amount").
		Joins("left join payments on payments.user_id = users.id").
		Where("DATE(payments.date) = ?", tomorrow).
		Count(&count).
		Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	body_message := "these are the payments that will be executed tomorrow. \r\n"
	if count > 0 {
		for rows.Next() {
			rows.Scan(&email, &name, &net_amount, &gross_amount)
			details := fmt.Sprintf("payment: %s \r\n net amount: %.2f \r\n gross amount: %.2f \r\n", name, net_amount, gross_amount)
			listPayments[email] += details
		}
		for key, val := range listPayments {
			mailsend.SendMail([]string{key}, "payment alert", body_message+val)
		}
	} else {
		log.Print("there are no payments")
	}
}
