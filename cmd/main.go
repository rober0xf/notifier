package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/rober0xf/notifier/cmd/api"
	_ "github.com/rober0xf/notifier/docs"
	"github.com/rober0xf/notifier/internal/infraestructure/scheduler"
)

// @title           Notifier API
// @version         1.0
// @description     Notifier API
// @host            localhost:3000
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	godotenv.Load()
	scheduler.InitCron()

	server, err := api.NewAPIServer(":3000")
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	log.Fatal(server.Run())
}
