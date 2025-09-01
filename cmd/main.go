package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rober0xf/notifier/internal/adapters/handlers/users"
	"github.com/rober0xf/notifier/internal/adapters/httpmethod"
	"github.com/rober0xf/notifier/internal/adapters/storage"
	database "github.com/rober0xf/notifier/internal/ports/db"
	cronjob "github.com/rober0xf/notifier/internal/scheduler"
	"github.com/rober0xf/notifier/internal/services/auth"
	userService "github.com/rober0xf/notifier/internal/services/users"
)

func main() {
	cronjob.InitCron()

	db, err := database.ConnectSQLite()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	// init repos
	userRepo := storage.NewUserRepository(db)

	// init services
	jwtKey := database.JwtKey
	authSvc := auth.NewAuthService(userRepo, jwtKey)
	userSvc := userService.NewUserService(userRepo, jwtKey)

	// init handlers
	userHandler := users.NewUserHandler(userSvc, authSvc)

	router := httpmethod.SetupRoutes(userHandler, jwtKey)

	fmt.Println("running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
