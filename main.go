package main

import (
	"fmt"
	"goapi/dbconnect"
	"goapi/handlers"
	"goapi/routes"
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

func main() {
	if err := LoadTemplates(); err != nil {
		log.Printf("error loading templates: %v", err)
		return
	}

	db, err := dbconnect.Connect()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	handlers.SendAlert()

	r := routes.InitRouter(db, tmpl)

	fmt.Println("running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func LoadTemplates() error {
	var err error
	tmpl, err = template.ParseGlob("templates/*.templ")

	if err != nil {
		return err
	}

	return nil
}
