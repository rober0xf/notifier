package main

// import (
// 	"database/sql"
// 	"log"
// 	"os"

// 	"github.com/golang-migrate/migrate/v4"
// 	"github.com/golang-migrate/migrate/v4/database/sqlite3"
// 	"github.com/golang-migrate/migrate/v4/source/file"
// 	_ "github.com/golang-migrate/migrate/v4/source/file"
// 	_ "github.com/mattn/go-sqlite3"
// )

// func main() {
// 	if len(os.Args) < 2 {
// 		log.Fatal("migration direction not provided. 'up' or 'down'")
// 	}

// 	direction := os.Args[1]

// 	db, err := sql.Open("sqlite3", "./data.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer db.Close()

// 	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", driver)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	switch direction {
// 	case "up":
// 		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
// 			log.Fatal(err)
// 		}
// 		log.Println("migrations applied successfully")
// 	case "down":
// 		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
// 			log.Fatal(err)
// 		}
// 		log.Println("migrations rolled back successfully")
// 	default:
// 		log.Fatal("invalid direction. use 'up' or 'down'")
// 	}
// }
