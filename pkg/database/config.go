package database

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rober0xf/notifier/pkg/email"
)

type PostgresConfig struct {
	DB_NAME string
	DB_USER string
	DB_PASS string
	DB_HOST string
	DB_PORT string
}

var JwtKey []byte
var MailSender email.EmailSender

func GetConfig() PostgresConfig {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading env file: %v", err)
	}

	var config PostgresConfig
	config.DB_NAME = GetEnvOrFatal("POSTGRES_DB")
	config.DB_USER = GetEnvOrFatal("POSTGRES_USER")
	config.DB_PASS = GetEnvOrFatal("POSTGRES_PASSWORD")
	config.DB_HOST = fallbackEnv("POSTGRES_HOST", "localhost")
	config.DB_PORT = fallbackEnv("POSTGRES_PORT", "5432")

	// validate that port is a number
	if _, err := strconv.Atoi(config.DB_PORT); err != nil {
		log.Fatalf("PORT env is not a number: %v", err)
	}

	jwt_key := strings.TrimSpace(GetEnvOrFatal("JWT_KEY"))
	JwtKey = []byte(jwt_key)

	return config
}

func GetEnvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("env variable %s empty", key)
	}
	return value
}

func fallbackEnv(key string, default_value string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return default_value
}
