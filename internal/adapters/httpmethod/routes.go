package httpmethod

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/authentication"
)

type UserHandler interface {
	// GET
	GetByEmailEmpty(c *gin.Context)
	GetByEmail(c *gin.Context)
	GetByID(c *gin.Context)
	GetAll(c *gin.Context)
	GetVerificationEmail(c *gin.Context)

	// POST
	Create(c *gin.Context)
	Login(c *gin.Context)

	// PUT-PATCH
	Update(c *gin.Context)

	// DELETE
	Delete(c *gin.Context)
}

type PaymentHandler interface {
	// GET
	GetAllPayments(c *gin.Context)
	GetPaymentByID(c *gin.Context)
	GetAllPaymentsFromUser(c *gin.Context)

	// POST
	CreatePayment(c *gin.Context)

	// PUT-PATCH
	UpdatePayment(c *gin.Context)

	// DELETE
	DeletePayment(c *gin.Context)
}

func SetupRoutes(userHandler UserHandler, paymentHandler PaymentHandler, jwtKey []byte) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           48 * time.Hour,
	}))

	v1 := r.Group("/v1")
	protected := v1.Group("/auth")

	// again to protected routes to ensure it runs before auth
	protected.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           48 * time.Hour,
	}))
	protected.Use(authentication.JWTMiddleware(jwtKey))

	setupUsersRoutes(v1, protected, userHandler)
	setupPaymentsRoutes(protected, paymentHandler)

	r.Static("/static", "./frontend/dist")

	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	return r
}

func setupUsersRoutes(v1, protected *gin.RouterGroup, userHandler UserHandler) {
	// TODO: make admin and self user handlers
	public := v1.Group("/users")
	auth := protected.Group("/users")

	/* ADMIN ROUTES */
	public.GET("", userHandler.GetAll)
	public.GET("/email_verification/:email/:hash", userHandler.GetVerificationEmail)
	public.POST("/register", userHandler.Create)
	public.POST("/login", userHandler.Login)

	auth.GET("/email", userHandler.GetByEmailEmpty)
	auth.GET("/email/:email", userHandler.GetByEmail)
	auth.GET("/:id", userHandler.GetByID)
	auth.PUT("/:id", userHandler.Update)
	auth.DELETE("/:id", userHandler.Delete)
}

func setupPaymentsRoutes(protected *gin.RouterGroup, paymentHandler PaymentHandler) {
	r := protected.Group("/payments")

	r.GET("/", paymentHandler.GetAllPayments)
	r.POST("", paymentHandler.CreatePayment)
	r.GET("/email", paymentHandler.GetAllPaymentsFromUser)
	r.GET("/:id", paymentHandler.GetPaymentByID)
	r.PUT("/:id", paymentHandler.UpdatePayment)
	r.DELETE("/:id", paymentHandler.DeletePayment)
}
