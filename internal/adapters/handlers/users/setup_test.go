package users

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rober0xf/notifier/internal/adapters/storage"
	"github.com/rober0xf/notifier/internal/adapters/testutils"
	"github.com/rober0xf/notifier/internal/ports"
	"github.com/rober0xf/notifier/internal/services/auth"
	"github.com/rober0xf/notifier/internal/services/users"

	"testing"
)

type TestDependencies struct {
	router      *gin.Engine
	db          *pgxpool.Pool
	userRepo    ports.UserRepository
	userService ports.UserService
	authService ports.AuthService
}

func SetupTestDependencies(t *testing.T) *TestDependencies {
	db := testutils.SetupTestDB(t)

	jwt := "test_secret"
	userRepo := storage.NewUserRepository(db)
	userService := users.NewUserService(userRepo, []byte(jwt))
	authService := auth.NewAuthService(userRepo, []byte(jwt))

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userHandler := NewUserHandler(userService, authService)

	// register routes
	router.POST("/users", userHandler.Create)
	router.POST("/users/login", userHandler.Login)
	router.GET("/users/:id", userHandler.GetByID)
	router.GET("/users/empty/:email", userHandler.GetByEmailEmpty)
	router.GET("/users/email/:email", userHandler.GetByEmail)
	router.GET("/users", userHandler.GetAll)
	router.PUT("/users/:id", userHandler.Update)
	router.DELETE("/users/:id", userHandler.Delete)

	return &TestDependencies{
		router:      router,
		db:          db,
		userRepo:    userRepo,
		userService: userService,
		authService: authService,
	}
}
