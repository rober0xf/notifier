package dbconnect

import (
	"fmt"
	"goapi/config"
	"goapi/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	config := config.GetConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB_USER,
		config.DB_PASS,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect to db: %v", err.Error())
	}

	fmt.Println("database connected")
	fmt.Println("migrating schema")

	DB.AutoMigrate(models.User{})
	DB.AutoMigrate(models.Category{})

	fmt.Println("schema migrated")

	return DB, nil
}
