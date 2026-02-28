package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	routes "github.com/rober0xf/notifier/internal/delivery/http"
	"github.com/rober0xf/notifier/internal/domain/repository"
	"github.com/rober0xf/notifier/internal/usecase/payment"
	"github.com/rober0xf/notifier/internal/usecase/user"
	"github.com/rober0xf/notifier/pkg/auth"
	"github.com/rober0xf/notifier/pkg/email"
)

func Serve(router *gin.Engine) error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", 3000),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("starting server on port: %d", 3000)

	return server.ListenAndServe()
}

func BuildUserRoutes(
	userRepo repository.UserRepository,
	tokenGen auth.TokenGenerator,
	emailSender email.EmailSender,
	disposableEmailChecker []string,
	baseURL string,
) *routes.UserHandler {

	createUserUC := user.NewCreateUserUseCase(userRepo, emailSender, disposableEmailChecker, baseURL)
	loginUC := user.NewLoginUseCase(userRepo, tokenGen)
	getUserByIDUC := user.NewGetUserByIDUseCase(userRepo)
	getUserByEmailUC := user.NewGetUserByEmailUseCase(userRepo)
	getAllUsersUC := user.NewGetAllUsersUseCase(userRepo)
	updateUserUC := user.NewUpdateUserUseCase(userRepo)
	deleteUserUC := user.NewDeleteUserUseCase(userRepo)
	verifyEmailUC := user.NewVerifyEmailUseCase(userRepo)

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
