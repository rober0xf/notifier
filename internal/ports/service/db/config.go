package database

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rober0xf/notifier/internal/models"
)

type Config struct {
	DB_NAME string
	DB_USER string
	DB_PASS string
	DB_HOST string
	DB_PORT string
}

var JwtKey []byte
var MailSender models.MailSender

func GetConfig() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading env file: %v", err)
	}

	var config Config
	config.DB_NAME = get_env_or_fatal("DB_NAME")
	config.DB_USER = get_env_or_fatal("DB_USER")
	config.DB_PASS = get_env_or_fatal("DB_PASSWORD")
	config.DB_HOST = fallback_env("DB_HOST", "localhost")
	config.DB_PORT = fallback_env("DB_PORT", "5432")

	// validate that port is a number
	if _, err := strconv.Atoi(config.DB_PORT); err != nil {
		log.Fatalf("PORT env is not a number: %v", err)
	}

	MailSender.Host = get_env_or_fatal("SMTP_HOST")
	MailSender.Port = get_env_or_fatal("SMTP_PORT")
	MailSender.Username = get_env_or_fatal("SMTP_USERNAME")
	MailSender.Password = get_env_or_fatal("SMTP_PASSWORD")

	// temporal key to check
	jwt_key := get_env_or_fatal("JWT_KEY")

	JwtKey = []byte(jwt_key)

	return config
}

func get_env_or_fatal(key string) string {
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
