package main

import (
	"fmt"
	"goapi/dbconnect"
	"goapi/routes"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome from homehandler"))
}

func main() {
	db, err := dbconnect.Connect()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	r := routes.InitRouter(db)

	r.HandleFunc("/home", homeHandler)

	fmt.Println("running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
