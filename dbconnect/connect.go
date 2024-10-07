package dbconnect

import (
	"fmt"
	"goapi/config"
	"goapi/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() (*gorm.DB, error) {
	config := config.GetConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB_USER,
		config.DB_PASS,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME)

	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect to db: %v", err.Error())
	}

	fmt.Println("database connected")
	fmt.Println("migrating schema")

	db.AutoMigrate(models.User{})

	fmt.Println("schema migrated")

	return db, nil
}
