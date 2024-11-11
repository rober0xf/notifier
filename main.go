package main

import (
	"fmt"
	"goapi/cronjob"
	"goapi/dbconnect"
	"goapi/routes"
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

func main() {
	go cronjob.InitCron()

	db, err := dbconnect.Connect()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	r := routes.InitRouter(db, tmpl)

	fmt.Println("running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
