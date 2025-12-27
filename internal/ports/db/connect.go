package database

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

var DB *pgxpool.Pool

func InitPostgres() (*pgxpool.Pool, error) {
	config := GetConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DB_USER,
		url.QueryEscape(config.DB_PASS),
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging: %w", err)
	}

	// migrations
	sqlDB := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("error closing db: %v", err)
		}
	}()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("error with driver: %w", err)
	}

	migrationsPath := GetEnvOrFatal("MIGRATIONS_PATH")
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("error while init migrate: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("error migrating: %w", err)
	}

	fmt.Println("database initialized")
	return pool, nil
}

// func ConnectSQLite() (*gorm.DB, error) {
// 	database_path, err := filepath.Abs("./database.db")
// 	if err != nil {
// 		return nil, fmt.Errorf("could not read database path: %v", err)
// 	}
//
// 	db, err := gorm.Open(sqlite.Open(database_path), &gorm.Config{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connecto to the database: %v", err)
// 	}
// 	fmt.Println("connected to sqlite")
//
// 	err = db.AutoMigrate(
// 		&domain.User{},
// 		&domain.Payment{},
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to migrate models: %v", err)
// 	}
//
// 	fmt.Println("schema migrated")
//
// 	return db, nil
// }
