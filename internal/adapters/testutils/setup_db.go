package testutils

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	database "github.com/rober0xf/notifier/internal/ports/db"
	"github.com/stretchr/testify/require"
)

func SetupTestDB(t *testing.T) *pgxpool.Pool {
	if err := loadRootEnv(); err != nil {
		log.Printf("warning: could not load env file: %v", err)
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
	_, _ = db.Exec(ctx, "DROP SCHEMA public CASCADE; CREATE SCHEMA public;")

	// then we create
	users, err := os.ReadFile("../../../../sql/schemas/users.sql")
	require.NoError(t, err)

	payments, err := os.ReadFile("../../../../sql/schemas/payments.sql")
	require.NoError(t, err)

	_, err = db.Exec(ctx, string(users))
	if err != nil {
		t.Fatalf("failed to create users table: %v\nSQL:\n%s", err, string(users))
	}
	require.NoError(t, err)

	_, err = db.Exec(ctx, string(payments))
	if err != nil {
		t.Fatalf("failed to create payments table: %v\nSQL:\n%s", err, string(payments))
	}
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

func InsertTestUser(ctx context.Context, db *pgxpool.Pool, email, username string) (int, error) {
	var id int
	err := db.QueryRow(ctx, `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`, username, email, "hashedpassword").Scan(&id)
	return id, err
}
