package database

import (
	"fmt"
	"log"
	"github.com/rober0xf/notifier/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	// from config.go
	config := GetConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.DB_HOST,
		config.DB_USER,
		config.DB_PASS,
		config.DB_NAME,
		config.DB_PORT)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err.Error())
	}

	fmt.Println("database connected")
	fmt.Println("migrating schema")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Payment{},
	)
	if err != nil {
		log.Fatalf("failed to migrate models: %v", err.Error())
	}

	fmt.Println("schema migrated")
	return DB, nil
}
