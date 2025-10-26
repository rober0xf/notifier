package database

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rober0xf/notifier/internal/services/mail"
)

type PostgresConfig struct {
	DB_NAME string
	DB_USER string
	DB_PASS string
	DB_HOST string
	DB_PORT string
}

var JwtKey []byte
var MailSender mail.MailSender

func GetConfig() PostgresConfig {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading env file: %v", err)
	}

	var config PostgresConfig
	config.DB_NAME = GetEnvOrFatal("DB_NAME")
	config.DB_USER = GetEnvOrFatal("DB_USER")
	config.DB_PASS = GetEnvOrFatal("DB_PASSWORD")
	config.DB_HOST = fallback_env("DB_HOST", "localhost")
	config.DB_PORT = fallback_env("DB_PORT", "5432")

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

func fallback_env(key string, default_value string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return default_value
}
