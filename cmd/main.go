package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/rober0xf/notifier/cmd/api"
	"github.com/rober0xf/notifier/internal/infraestructure/scheduler"
)

func main() {
	godotenv.Load()
	scheduler.InitCron()

	server, err := api.NewAPIServer(":3000")
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	log.Fatal(server.Run())
}
