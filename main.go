package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome from homehandler"))
}

func main() {
	fmt.Println("hello from first line")

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)

	fmt.Println("running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
