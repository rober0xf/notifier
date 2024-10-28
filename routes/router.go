package routes

import (
	"goapi/config"
	"goapi/handlers"
	"goapi/middlewares"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, tmpl *template.Template) *mux.Router {
	r := mux.NewRouter()

	// subrouters
	userRouter := r.PathPrefix("/users").Subrouter()
	protectedRouter := r.PathPrefix("/protected").Subrouter()
	categoryRouter := r.PathPrefix("/categories").Subrouter()
	paymentRouter := r.PathPrefix("/payment").Subrouter()

	protectedRouter.Use(middlewares.JWTMiddleware) // apply the middleware
	categoryRouter.Use(middlewares.JWTMiddleware)
	paymentRouter.Use(middlewares.JWTMiddleware)

	r.HandleFunc("/", handlers.HomeHandler).Methods(http.MethodGet)
	r.HandleFunc("/login", handlers.LoginTemplate()).Methods(http.MethodGet)
	r.HandleFunc("/login", handlers.LoginTemplate()).Methods(http.MethodPost)

	// users
	userRouter.HandleFunc("", handlers.GetUser).Methods(http.MethodGet)
	userRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUser(db, w, r)
	}).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/users/{id}", handlers.GetUser).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateUser(db, w, r)
	}).Methods(http.MethodPut)
	protectedRouter.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteUser(db, w, r)
	}).Methods(http.MethodDelete)

	// categories
	categoryRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateCategory(db, w, r)
	}).Methods(http.MethodPost)
	categoryRouter.HandleFunc("", handlers.GetCategories).Methods(http.MethodGet)
	categoryRouter.HandleFunc("/{id}", handlers.GetCategories).Methods(http.MethodGet)
	categoryRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateCategory(db, w, r)
	}).Methods(http.MethodPut)
	categoryRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteCategory(db, w, r)
	}).Methods(http.MethodDelete)

	// payment
	paymentRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreatePayment(db, w, r)
	}).Methods(http.MethodPost)
	paymentRouter.HandleFunc("", handlers.GetPayment).Methods(http.MethodGet)
	paymentRouter.HandleFunc("/{id}", handlers.GetPayment).Methods(http.MethodGet)
	paymentRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdatePayment(db, w, r)
	}).Methods(http.MethodPut)
	paymentRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeletePayment(db, w, r)
	}).Methods(http.MethodDelete)

	protectedRouter.HandleFunc("/email", func(w http.ResponseWriter, r *http.Request) {
		config.TestMail(w, r)
	}).Methods(http.MethodGet)

	return r
}
