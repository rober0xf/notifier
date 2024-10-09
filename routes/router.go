package routes

import (
	"goapi/handlers"
	"goapi/middlewares"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, tmpl *template.Template) *mux.Router {
	r := mux.NewRouter()

	userRouter := r.PathPrefix("/users").Subrouter()
	protectedRouter := r.PathPrefix("/protected").Subrouter()
	protectedRouter.Use(middlewares.JWTMiddleware)
	categoryRouter := r.PathPrefix("/categories").Subrouter()

	r.HandleFunc("/login", handlers.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/", handlers.HomeHandler(tmpl)).Methods(http.MethodGet)

	// GET public routes
	userRouter.HandleFunc("", handlers.GetUser).Methods(http.MethodGet)
	categoryRouter.HandleFunc("", handlers.GetCategories).Methods(http.MethodGet)
	categoryRouter.HandleFunc("/{id}", handlers.GetCategories).Methods(http.MethodGet)

	// POST routes
	userRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUser(db, w, r)
	}).Methods(http.MethodPost)

	// GET protected routes
	protectedRouter.HandleFunc("/user/{id}", handlers.GetUser).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/test", middlewares.ProtectedTest).Methods(http.MethodGet)

	// PUT routes
	protectedRouter.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateUser(db, w, r)
	}).Methods(http.MethodPut)

	// DELETE routes
	protectedRouter.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteUser(db, w, r)
	}).Methods(http.MethodDelete)

	return r
}
