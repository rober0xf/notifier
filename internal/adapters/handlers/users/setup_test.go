package users

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rober0xf/notifier/internal/adapters/storage"
	"github.com/rober0xf/notifier/internal/ports"
	database "github.com/rober0xf/notifier/internal/ports/db"
	"github.com/rober0xf/notifier/internal/services/auth"
	"github.com/rober0xf/notifier/internal/services/users"
	"github.com/stretchr/testify/require"
)

type TestDependencies struct {
	router      *gin.Engine
	db          *pgxpool.Pool
	userRepo    ports.UserRepository
	userService ports.UserService
	authService ports.AuthService
}

func SetupTestDependencies(t *testing.T) *TestDependencies {
	db := SetupTestDB(t)

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

func SetupTestDB(t *testing.T) *pgxpool.Pool {
	if err := loadRootEnv(); err != nil {
		log.Fatalf("error loading env file: %v", err)
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		database.GetEnvOrFatal("POSTGRES_USER_TEST"),
		url.QueryEscape(database.GetEnvOrFatal("POSTGRES_PASSWORD_TEST")),
		database.GetEnvOrFatal("POSTGRES_HOST_TEST"),
		database.GetEnvOrFatal("POSTGRES_PORT_TEST"),
		database.GetEnvOrFatal("POSTGRES_NAME_TEST"),
	)

	db, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)

	err = db.Ping(context.Background())
	require.NoError(t, err)

	runMigrations(t, db)
	cleanDatabase(t, db)
	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func runMigrations(t *testing.T, db *pgxpool.Pool) {
	ctx := context.Background()

	// first we drop if exists
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS payments CASCADE")
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS users CASCADE")
	_, _ = db.Exec(ctx, "DROP TYPE IF EXISTS transaction_type CASCADE")
	_, _ = db.Exec(ctx, "DROP TYPE IF EXISTS category_type CASCADE")
	_, _ = db.Exec(ctx, "DROP TYPE IF EXISTS frequency_type CASCADE")

	// then we create
	users, err := os.ReadFile("../../../../sql/schemas/users.sql")
	require.NoError(t, err)

	payments, err := os.ReadFile("../../../../sql/schemas/payments.sql")
	require.NoError(t, err)

	_, err = db.Exec(ctx, string(users))
	require.NoError(t, err)

	_, err = db.Exec(ctx, string(payments))
	require.NoError(t, err)
}

func cleanDatabase(t *testing.T, db *pgxpool.Pool) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE users CASCADE")
	require.NoError(t, err)

	_, err = db.Exec(ctx, "TRUNCATE payments CASCADE")
	require.NoError(t, err)
}

func loadRootEnv() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			log.Printf("loading .env from: %s", envPath)
			return godotenv.Load(envPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return fmt.Errorf(".env not found in any parent directory")
		}
		dir = parent
	}
}
