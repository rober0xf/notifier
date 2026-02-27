package integration

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
	"github.com/rober0xf/notifier/pkg/database"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	if err := loadRootEnv(); err != nil {
		log.Printf("warning: could not load env file: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		database.GetEnvOrFatal("POSTGRES_USER_TEST"),
		url.QueryEscape(database.GetEnvOrFatal("POSTGRES_PASSWORD_TEST")),
		database.GetEnvOrFatal("POSTGRES_HOST_TEST"),
		database.GetEnvOrFatal("POSTGRES_PORT_TEST"),
		database.GetEnvOrFatal("POSTGRES_DB_TEST"),
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

func insertTestUser(ctx context.Context, db *pgxpool.Pool, email string, username string) (int, error) {
	var id int
	err := db.QueryRow(ctx, `
		INSERT INTO users (username, email, password, active)
		VALUES ($1, $2, $3)
		RETURNING id
	`, username, email, "hashedpassword", true).Scan(&id)
	return id, err
}

func insertTestPayment(ctx context.Context, db *pgxpool.Pool, userID int, name string, amount float64) (int, error) {
	var id int
	err := db.QueryRow(ctx, `
		INSERT INTO payments (user_id, name, amount, type, category, date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, userID, name, amount, "expense", "electronics", "2026-01-01").Scan(&id)
	return id, err
}

func runMigrations(t *testing.T, db *pgxpool.Pool) {
	ctx := context.Background()

	_, _ = db.Exec(ctx, `DROP TABLE IF EXISTS payments CASCADE;`)
	_, _ = db.Exec(ctx, `DROP TABLE IF EXISTS users CASCADE;`)
	_, _ = db.Exec(ctx, `DROP TYPE IF EXISTS transaction_type CASCADE;`)
	_, _ = db.Exec(ctx, `DROP TYPE IF EXISTS category_type CASCADE;`)
	_, _ = db.Exec(ctx, `DROP TYPE IF EXISTS frequency_type CASCADE;`)

	users, err := os.ReadFile("../../../../sql/schemas/users.sql")
	require.NoError(t, err)

	_, err = db.Exec(ctx, string(users))
	require.NoError(t, err)

	payments, err := os.ReadFile("../../../../sql/schemas/payments.sql")
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
