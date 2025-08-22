package main

import (
	"fmt"
	"log"
	"net/http"

	userHandler "github.com/rober0xf/notifier/cmd/api/handlers/users"
	database "github.com/rober0xf/notifier/internal/ports/db"
	"github.com/rober0xf/notifier/internal/routes"
	cronjob "github.com/rober0xf/notifier/internal/scheduler"
	"github.com/rober0xf/notifier/internal/services/users"
)

func main() {
	cronjob.InitCron()

	db, err := database.ConnectSQLite()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	userService := &users.Service{}
	userHandler := userHandler.Handler{
		UserService: userService,
	}

	r := routes.InitRouter(db)

	fmt.Println("running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
