package integration

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rober0xf/notifier/internal/delivery/http"
	"github.com/rober0xf/notifier/internal/domain/repository"
	"github.com/rober0xf/notifier/internal/infraestructure/persistance/postgres"
	"github.com/rober0xf/notifier/internal/usecase/payment"
	"github.com/rober0xf/notifier/internal/usecase/user"
	"github.com/rober0xf/notifier/pkg/auth"
	"github.com/rober0xf/notifier/pkg/email"

	"testing"
)

type TestUserDependencies struct {
	router   *gin.Engine
	db       *pgxpool.Pool
	userRepo repository.UserRepository
}

type TestPaymentDependencies struct {
	router      *gin.Engine
	db          *pgxpool.Pool
	paymentRepo repository.PaymentRepository
	userRepo    payment.UserIDGetter
}

func setupTestDependencies(t *testing.T) (*TestUserDependencies, *TestPaymentDependencies) {
	t.Helper()

	db := setupTestDB(t)

	userRepo := postgres.NewUserRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)

	jwtKey := []byte("test-jwt")
	tokenGen := auth.NewJWTGenerator(jwtKey, 160)
	authMiddleware := auth.AuthMiddleware(tokenGen, auth.SessionCookieName)

	emailSender := email.NewMockSender()

	createUserUC := user.NewCreateUserUseCase(userRepo, emailSender, []string{}, "http://localhost:8080")
	loginUserUC := user.NewLoginUseCase(userRepo, tokenGen)
	getUserByIDUC := user.NewGetUserByIDUseCase(userRepo)
	getUserByEmailUC := user.NewGetUserByEmailUseCase(userRepo)
	getAllUsersUC := user.NewGetAllUsersUseCase(userRepo)
	updateUserUC := user.NewUpdateUserUseCase(userRepo)
	verifyEmailUC := user.NewVerifyEmailUseCase(userRepo)
	deleteUserUC := user.NewDeleteUserUseCase(userRepo)

	userHandler := http.NewUserHandler(
		createUserUC,
		loginUserUC,
		getUserByIDUC,
		getUserByEmailUC,
		getAllUsersUC,
		updateUserUC,
		deleteUserUC,
		verifyEmailUC,
		tokenGen,
	)

	createPaymentUC := payment.NewCreatePaymentUseCase(paymentRepo)
	getPaymentByIDUC := payment.NewGetPaymentByIDUseCase(paymentRepo)
	getAllPaymentsFromUserUC := payment.NewGetAllPaymentsFromUserUseCase(paymentRepo, userRepo)
	getAllPaymentsUC := payment.NewGetAllPaymentsUseCase(paymentRepo)
	updatePaymentUC := payment.NewUpdatePaymentUseCase(paymentRepo)
	deletePaymentUC := payment.NewDeletePaymentUseCase(paymentRepo)

	paymentHandler := http.NewPaymentHandler(
		createPaymentUC,
		getPaymentByIDUC,
		getAllPaymentsFromUserUC,
		getAllPaymentsUC,
		updatePaymentUC,
		deletePaymentUC,
	)

	gin.SetMode(gin.TestMode)
	router := http.SetupRoutes(userHandler, paymentHandler, authMiddleware)

	return &TestUserDependencies{
			router:   router,
			db:       db,
			userRepo: userRepo,
		}, &TestPaymentDependencies{
			router:      router,
			db:          db,
			paymentRepo: paymentRepo,
			userRepo:    userRepo,
		}
}
