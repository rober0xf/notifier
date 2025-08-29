package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/handlers/users"
	"github.com/rober0xf/notifier/internal/adapters/httpmethod"
	"github.com/rober0xf/notifier/internal/adapters/storage"
	database "github.com/rober0xf/notifier/internal/ports/db"
	cronjob "github.com/rober0xf/notifier/internal/scheduler"
	userService "github.com/rober0xf/notifier/internal/services/users"
)

func main() {
	cronjob.InitCron()

	db, err := database.ConnectSQLite()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	userSvc := &userService.Service{
		Repo: storage.NewUserRepository(db),
	}
	userHandler := users.NewUserHandler(*userSvc)

	router := gin.Default()
	jwtKey := []byte("your-secret-key")

	httpmethod.SetupRoutes(userHandler, jwtKey)

	fmt.Println("running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
