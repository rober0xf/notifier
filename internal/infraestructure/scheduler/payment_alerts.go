package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rober0xf/notifier/pkg/database"
	"github.com/rober0xf/notifier/pkg/email"
)

var ms email.SMTPSender

func init() {
	ms = *email.NewSMTPSender(
		database.GetEnvOrFatal("SMTP_HOST"),
		database.GetEnvOrFatal("SMTP_PORT"),
		database.GetEnvOrFatal("SMTP_USERNAME"),
		database.GetEnvOrFatal("SMTP_PASSWORD"),
	)
}

type payment struct {
	Email   string
	Name    string
	Amount  float64
	PayType string
}

func InitCron() {
	go func() {
		s, err := gocron.NewScheduler(gocron.WithLocation(time.Local))
		if err != nil {
			log.Fatalf("error creating scheduler: %v", err)
			return
		}

		// set daily cronjob
		_, err = s.NewJob(
			gocron.DailyJob(1, gocron.NewAtTimes(
				gocron.NewAtTime(07, 00, 00),
			),
			),
			gocron.NewTask(SendDailyAlert),
		)
		if err != nil {
			log.Printf("error creating daily job: %v", err)
			return
		}
		log.Println("daily cronjob set")

		// set weekly cronjob
		_, err = s.NewJob(
			gocron.WeeklyJob(1, gocron.NewWeekdays(
				time.Monday),
				gocron.NewAtTimes(gocron.NewAtTime(07, 00, 00))),
			gocron.NewTask(SendWeeklyAlert))
		if err != nil {
			log.Printf("error creating weekly job: %v", err)
			return
		}
		log.Println("weekly cronjob set")

		// set monthly cronjob
		_, err = s.NewJob(
			gocron.MonthlyJob(1, gocron.NewDaysOfTheMonth(1),
				gocron.NewAtTimes(gocron.NewAtTime(07, 00, 00))),
			gocron.NewTask(SendMonthlyAlert))
		if err != nil {
			log.Printf("error creating monthly job: %v", err)
			return
		}
		log.Println("monthly cronjob set")

		s.Start()

		select {}
	}()
}

func SendPaymentAlert(title string, target_date time.Time) error {
	ctx := context.Background()
	log.Printf("starting SendPaymentAlert for date: %s", target_date.Format("2006-01-02"))

	query := `
			SELECT users.email, payments.name, payments.amount, payments.type
			FROM users
			LEFT JOIN payments ON payments.user_id = users.id
			WHERE DATE(payments.due_date) = DATE($1)
			`
	rows, err := database.DB.Query(ctx, query, target_date)
	if err != nil {
		log.Printf("database error: %v", err)
		return err
	}
	defer rows.Close()

	list_payments := make(map[string]string)
	has_payments := false
	payment_count := 0

	for rows.Next() {
		var p payment
		if err := rows.Scan(&p.Email, &p.Name, &p.Amount, &p.PayType); err != nil {
			log.Printf("error scanning row: %s", err)
			continue
		}

		has_payments = true
		payment_count++
		details := fmt.Sprintf("Payment: %s\r\nAmount: %.2f\r\nType: %s\r\n", p.Name, p.Amount, p.PayType)
		list_payments[p.Email] += details
		log.Printf("found payment: %s for %s (%f)", p.Name, p.Email, p.Amount)
	}

	if !has_payments {
		log.Printf("there are no payments due on %s", target_date.Format("2006-01-02"))
		return nil
	}

	bodyMessage := fmt.Sprintf("These are the payments that will be executed on %s:\r\n", target_date.Format("2006-01-02"))
	emails_sent := 0
	emails_failed := 0

	for email, val := range list_payments {
		log.Printf("trying to send email to: %s", email)
		if err := ms.Send([]string{email}, title, bodyMessage+val); err != nil {
			log.Printf("error sending email to %s: %v", email, err)
			emails_failed++
		} else {
			log.Printf("email sent to: %s", email)
			emails_sent++
		}
	}

	log.Printf("email stats - Sent: %d, Failed: %d", emails_sent, emails_failed)
	return nil
}

func SendDailyAlert() {
	log.Println("========== SendDailyAlert triggered ==========")

	tomorrow := time.Now().AddDate(0, 0, 1)
	if err := SendPaymentAlert("daily payment alert", tomorrow); err != nil {
		log.Printf("SendDailyAlert error: %v", err)
	} else {
		log.Println("========== SendDailyAlert finished ==========")
	}
}

func SendWeeklyAlert() {
	log.Println("========== SendWeeklyAlert triggered ==========")

	next_week := time.Now().AddDate(0, 0, 7)
	if err := SendPaymentAlert("weekly payment alert", next_week); err != nil {
		log.Printf("SendWeeklyAlert error: %v", err)
	} else {
		log.Println("========== SendWeeklyAlert finished ==========")
	}
}

func SendMonthlyAlert() {
	log.Println("========== SendMonthlyAlert triggered ==========")

	next_month := time.Now().AddDate(0, 1, 0)
	if err := SendPaymentAlert("monthly payment alert", next_month); err != nil {
		log.Printf("SendMonthlyAlert error: %v", err)
	} else {
		log.Println("========== SendMonthlyAlert finished ==========")
	}
}
