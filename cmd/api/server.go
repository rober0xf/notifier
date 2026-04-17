package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	routes "github.com/rober0xf/notifier/internal/delivery/http"
	"github.com/rober0xf/notifier/internal/domain/repository"
	"github.com/rober0xf/notifier/internal/infraestructure/persistance/postgres"
	"github.com/rober0xf/notifier/internal/usecase/payment"
	"github.com/rober0xf/notifier/internal/usecase/user"
	"github.com/rober0xf/notifier/pkg/auth"
	"github.com/rober0xf/notifier/pkg/database"
	"github.com/rober0xf/notifier/pkg/email"
)

type APIServer struct {
	addr   string
	router http.Handler
}

func NewAPIServer(addr string) (*APIServer, error) {
	_ = database.GetConfig()

	db, err := database.InitPostgres()
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}

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

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is not set")
	}

	// handlers
	userHandler := BuildUserRoutes(userRepo, tokenGen, emailSender, disposableEmails, baseURL, googleClientID)
	paymentHandler := BuildPaymentRoutes(paymentRepo, userRepo)

	authMiddleware := auth.AuthMiddleware(tokenGen, "access_token")

	router := routes.SetupRoutes(userHandler, paymentHandler, authMiddleware)

	return &APIServer{
		addr:   addr,
		router: router,
	}, nil
}

func (s *APIServer) Run() error {
	server := &http.Server{
		Addr:         s.addr,
		Handler:      s.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Printf("server running on %s", s.addr)

	return server.ListenAndServe()
}

func BuildUserRoutes(
	userRepo repository.UserRepository,
	tokenGen auth.TokenGenerator,
	emailSender email.EmailSender,
	disposableEmailChecker []string,
	baseURL string,
	googleClientID string,
) *routes.UserHandler {

	createUserUC := user.NewCreateUserUseCase(userRepo, emailSender, disposableEmailChecker, baseURL)
	loginUC := user.NewLoginUseCase(userRepo, tokenGen)
	getUserByIDUC := user.NewGetUserByIDUseCase(userRepo)
	getUserByEmailUC := user.NewGetUserByEmailUseCase(userRepo)
	getAllUsersUC := user.NewGetAllUsersUseCase(userRepo)
	updateUserUC := user.NewUpdateUserUseCase(userRepo)
	deleteUserUC := user.NewDeleteUserUseCase(userRepo)
	verifyEmailUC := user.NewVerifyEmailUseCase(userRepo)
	oauthUC := user.NewOAuthUseCase(userRepo)
	googleVerifier := auth.NewGoogleVerifier(googleClientID)

	return routes.NewUserHandler(
		createUserUC,
		loginUC,
		getUserByIDUC,
		getUserByEmailUC,
		getAllUsersUC,
		updateUserUC,
		deleteUserUC,
		verifyEmailUC,
		tokenGen,
		oauthUC,
		googleVerifier,
	)
}

func BuildPaymentRoutes(
	paymentRepo repository.PaymentRepository,
	userRepo payment.UserIDGetter,
) *routes.PaymentHandler {

	createPaymentUC := payment.NewCreatePaymentUseCase(paymentRepo)
	getAllPaymentsUC := payment.NewGetAllPaymentsUseCase(paymentRepo)
	getPaymentByIDUC := payment.NewGetPaymentByIDUseCase(paymentRepo)
	getAllPaymentsFromUserUC := payment.NewGetAllPaymentsFromUserUseCase(paymentRepo, userRepo)
	updatePaymentUC := payment.NewUpdatePaymentUseCase(paymentRepo)
	deletePaymentUC := payment.NewDeletePaymentUseCase(paymentRepo)

	return routes.NewPaymentHandler(
		createPaymentUC,
		getPaymentByIDUC,
		getAllPaymentsFromUserUC,
		getAllPaymentsUC,
		updatePaymentUC,
		deletePaymentUC,
	)
}
