package payments

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rober0xf/notifier/internal/adapters/storage"
	"github.com/rober0xf/notifier/internal/adapters/testutils"
	"github.com/rober0xf/notifier/internal/ports"
	"github.com/rober0xf/notifier/internal/services/auth"
	"github.com/rober0xf/notifier/internal/services/payments"
)

type TestDependencies struct {
	router         *gin.Engine
	db             *pgxpool.Pool
	paymentRepo    ports.PaymentRepository
	paymentService ports.PaymentService
	authService    ports.AuthService
}

func SetupTestDependencies(t *testing.T) *TestDependencies {
	db := testutils.SetupTestDB(t)

	jwt := "test_secret"
	paymentRepo := storage.NewPaymentRepository(db)
	paymentSrv := payments.NewPayments(paymentRepo)
	authService := auth.NewAuthService(paymentRepo, []byte(jwt))

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Next()
	})

	paymentHandler := NewPaymentHandler(paymentSrv, authService)

	// register routes
	router.POST("/payments", paymentHandler.CreatePayment)
	router.GET("/payments/email", paymentHandler.GetAllPaymentsFromUser)
	router.GET("/payments/:id", paymentHandler.GetPaymentByID)
	router.GET("/payments", paymentHandler.GetAllPayments)
	router.PUT("/payments/:id", paymentHandler.UpdatePayment)
	router.DELETE("/payments/:id", paymentHandler.DeletePayment)

	return &TestDependencies{
		router:         router,
		db:             db,
		paymentRepo:    paymentRepo,
		paymentService: paymentSrv,
		authService:    authService,
	}
}
