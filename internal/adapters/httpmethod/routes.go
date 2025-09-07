package httpmethod

import (
	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
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
	UpdatePayment(c *gin.Context)

	// DELETE
	DeletePayment(c *gin.Context)
}

func SetupRoutes(userHandler UserHandler, paymentHandler PaymentHandler, jwtKey []byte) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	protected := v1.Group("/auth")
	protected.Use(authentication.JWTMiddleware(jwtKey))

	setupUsersRoutes(v1, protected, userHandler)
	setupPaymentsRoutes(protected, paymentHandler)

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
			httphelpers.EmailParameterNotProvided(ctx)
			return
		}
		userHandler.GetUser(email, ctx)
	})

	// get user by id
	auth.GET("/:id", func(ctx *gin.Context) {
		id := ctx.Query("id")
		if id == "" {
			httphelpers.IDParameterNotProvided(ctx)
			return
		}
		userHandler.GetUserByID(id, ctx)
	})

	// update user
	auth.PUT("/:id", func(ctx *gin.Context) {
		id := ctx.Query("id")
		if id == "" {
			httphelpers.IDParameterNotProvided(ctx)
			return
		}
		userHandler.UpdateUser(ctx)
	})

	// delete user
	auth.DELETE("/:id", func(ctx *gin.Context) {
		id := ctx.Query("id")
		if id == "" {
			httphelpers.IDParameterNotProvided(ctx)
			return
		}
		userHandler.DeleteUser(ctx)
	})
}

func setupPaymentsRoutes(protected *gin.RouterGroup, paymentHandler PaymentHandler) {
	r := protected.Group("/payments")

	// get method
	r.GET("/", paymentHandler.GetAllPayments)
	r.GET("/:id", func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			httphelpers.IDParameterNotProvided(c)
			return
		}
		paymentHandler.GetPaymentByID(c)
	})

	// post method
	r.POST("/", paymentHandler.CreatePayment)

	// put method
	r.PUT("/:id", func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			httphelpers.IDParameterNotProvided(c)
			return
		}
		paymentHandler.UpdatePayment(c)
	})

	// delete method
	r.DELETE("/:id", func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			httphelpers.IDParameterNotProvided(c)
			return
		}
		paymentHandler.DeletePayment(c)
	})
}
