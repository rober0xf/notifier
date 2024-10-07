package routes

import (
	"goapi/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()
	userRouter := r.PathPrefix("/users").Subrouter()

	// GET public routes
	userRouter.HandleFunc("", handlers.GetUser).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id}", handlers.GetUser).Methods(http.MethodGet)

	// POST routes
	userRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUser(db, w, r)
	}).Methods(http.MethodPost)

	// PUT routes
	userRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateUser(db, w, r)
	}).Methods(http.MethodPut)

	userRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteUser(db, w, r)
	}).Methods(http.MethodDelete)

	return r
}
