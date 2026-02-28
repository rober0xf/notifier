package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rober0xf/notifier/cmd/api"
	routes "github.com/rober0xf/notifier/internal/delivery/http"
	"github.com/rober0xf/notifier/internal/infraestructure/persistance/postgres"
	"github.com/rober0xf/notifier/internal/infraestructure/scheduler"
	"github.com/rober0xf/notifier/pkg/auth"
	"github.com/rober0xf/notifier/pkg/database"
	"github.com/rober0xf/notifier/pkg/email"
)

func main() {
	_ = database.GetConfig()
	db, err := database.InitPostgres()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	scheduler.InitCron()

	// init repos
	userRepo := postgres.NewUserRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)

	// infra
	jwtKey := database.JwtKey
	tokenGen := auth.NewJWTGenerator(jwtKey, 24)
	emailSender := email.NewSMTPSender(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)
	disposableEmails := email.MustDisposableEmail()
	baseURL := os.Getenv("BASE_URL")

	// handlers
	userHandler := api.BuildUserRoutes(userRepo, tokenGen, emailSender, disposableEmails, baseURL)
	paymentHandler := api.BuildPaymentRoutes(paymentRepo, userRepo)

	authMiddleware := auth.AuthMiddleware(tokenGen, "access_token")

	router := routes.SetupRoutes(userHandler, paymentHandler, authMiddleware)

	fmt.Println("running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", router))
}
