package httpmethod

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	// GET
	GetUser(email string, c *gin.Context)
	GetAllUsers(c *gin.Context)
	GetUserByID(id string, c *gin.Context)

	// POST
	CreateUser(c *gin.Context)

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

	v1 := r.Group("/api/v1")
	protected := v1.Group("/auth")
	protected.Use()

	setupUsersRoutes(v1, protected, userHandler)

	return r
}

func setupUsersRoutes(v1, protected *gin.RouterGroup, userHandler UserHandler) {
	publicUsers := v1.Group("/users")

	// public routes
	publicUsers.POST("/", userHandler.CreateUser)
	publicUsers.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "login not implemented"})
	})

	// protected routes
	protectedUsers := protected.Group("/users")
	protectedUsers.GET("/", userHandler.GetAllUsers)
	protectedUsers.GET("/email", func(ctx *gin.Context) {
		email := ctx.Query("email")
		userHandler.GetUser(email, ctx)
	})
	protectedUsers.GET("/:id", func(ctx *gin.Context) {
		id := ctx.Query("id")
		userHandler.GetUserByID(id, ctx)
	})
	protectedUsers.PUT("/:id", func(c *gin.Context) {
		userHandler.UpdateUser(c)
	})
	protectedUsers.DELETE("/:id", userHandler.DeleteUser)
}

func setupPaymentsRoutes(v1, protected *gin.RouterGroup, paymentHandler UserHandler)
