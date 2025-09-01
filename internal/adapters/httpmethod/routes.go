package httpmethod

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/authentication"
)

type UserHandler interface {
	// GET
	GetUser(email string, c *gin.Context)
	GetAllUsers(c *gin.Context)
	GetUserByID(id string, c *gin.Context)

	// POST
	CreateUser(c *gin.Context)
	Login(c *gin.Context)

	// PUT-PATCH
	UpdateUser(c *gin.Context)

	// DELETE
	DeleteUser(c *gin.Context)
}

type PaymentHandler interface {
	// GET
	GetAllPayments(c *gin.Context)
	GetPaymentByID(c *gin.Context)

	// POST
	CreatePayment(c *gin.Context)

	// PUT-PATCH
	UpdateUser(c *gin.Context)

	// DELETE
	DeleteUser(c *gin.Context)
}

func SetupRoutes(userHandler UserHandler, jwtKey []byte) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	protected := v1.Group("/auth")
	protected.Use(authentication.JWTMiddleware(jwtKey))

	setupUsersRoutes(v1, protected, userHandler)

	return r
}

func setupUsersRoutes(v1, protected *gin.RouterGroup, userHandler UserHandler) {
	public := v1.Group("/users")
	auth := protected.Group("/users")

	// public routes
	public.POST("/register", userHandler.CreateUser)
	public.GET("", userHandler.GetAllUsers)
	public.POST("/login", userHandler.Login)

	// protected routes
	// get user by email
	auth.GET("/email", func(ctx *gin.Context) {
		email := ctx.Query("email")
		if email == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "email query parameter required"})
			return
		}
		userHandler.GetUser(email, ctx)
	})

	// get user by id
	auth.GET("/:id", func(ctx *gin.Context) {
		id := ctx.Query("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter required"})
			return
		}
		userHandler.GetUserByID(id, ctx)
	})

	// update user
	auth.PUT("/:id", func(ctx *gin.Context) {
		id := ctx.Query("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter required"})
			return
		}
		userHandler.UpdateUser(ctx)
	})
	auth.DELETE("/:id", userHandler.DeleteUser)
}

func setupPaymentsRoutes(v1, protected *gin.RouterGroup, paymentHandler UserHandler) {}
