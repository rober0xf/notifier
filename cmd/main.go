package main

import (
	"fmt"
	"log"
	"net/http"

	database "github.com/rober0xf/notifier/internal/db"
	"github.com/rober0xf/notifier/internal/routes"
	cronjob "github.com/rober0xf/notifier/internal/scheduler"
)

func main() {
	cronjob.InitCron()

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	r := routes.InitRouter(db)

	fmt.Println("running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
