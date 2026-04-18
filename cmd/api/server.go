package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	routes "github.com/rober0xf/notifier/internal/delivery/http"
	"github.com/rober0xf/notifier/internal/infraestructure/persistance/postgres"
	"github.com/rober0xf/notifier/pkg/auth"
	"github.com/rober0xf/notifier/pkg/database"
	"github.com/rober0xf/notifier/pkg/email"
)

type APIServer struct {
	addr        string
	router      http.Handler
	cancelClean context.CancelFunc
}

func NewAPIServer(addr string) (*APIServer, error) {
	db, err := database.InitPostgres()
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is not set")
	}

	// init repos
	userRepo := postgres.NewUserRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)

	tokenGen := auth.NewJWTGenerator(database.JwtKey, 24)
	emailSender := email.NewSMTPSender(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)

	// handlers
	userHandler := buildUserRoutes(userRepo, tokenGen, emailSender, googleClientID)
	paymentHandler := buildPaymentRoutes(paymentRepo, userRepo)
	router := routes.SetupRoutes(userHandler, paymentHandler, auth.AuthMiddleware(tokenGen, auth.SessionCookieName))

	// clean tokens
	ctx, cancel := context.WithCancel(context.Background())
	database.StartTokenCleanJob(ctx, userRepo)

	return &APIServer{addr: addr, router: router, cancelClean: cancel}, nil
}

func (s *APIServer) Run() error {
	server := &http.Server{
		Addr:         s.addr,
		Handler:      s.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("shutting down server")
		s.cancelClean()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	log.Printf("server running on %s", s.addr)
	return server.ListenAndServe()
}
