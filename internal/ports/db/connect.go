package database

import (
	"fmt"
	"path/filepath"

	"github.com/rober0xf/notifier/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// func ConnectPostgres() (*gorm.DB, error) {
// 	config := GetConfig()

// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
// 		config.DB_HOST,
// 		config.DB_USER,
// 		config.DB_PASS,
// 		config.DB_NAME,
// 		config.DB_PORT)

// 	var err error
// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("failed to connect to db: %v", err.Error())
// 	}

// 	fmt.Println("database connected")
// 	fmt.Println("migrating schema")

// 	err = DB.AutoMigrate(
// 		&domain.User{},
// 		&domain.Category{},
// 		&domain.Payment{},
// 	)
// 	if err != nil {
// 		log.Fatalf("failed to migrate models: %v", err.Error())
// 	}

// 	fmt.Println("schema migrated")
// 	return DB, nil
// }

func ConnectSQLite() (*gorm.DB, error) {
	database_path, err := filepath.Abs("./database.db")
	if err != nil {
		return nil, fmt.Errorf("could not read database path: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(database_path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connecto to the database: %v", err)
	}
	fmt.Println("connected to sqlite")

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Payment{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate models: %v", err)
	}

	fmt.Println("schema migrated")

	return db, nil
}
