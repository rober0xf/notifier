package routes

import (
	"github.com/gorilla/mux"
	"goapi/handlers"
	"goapi/middlewares"
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

func InitRouter(db *gorm.DB, tmpl *template.Template) *mux.Router {
	r := mux.NewRouter()
	store := &handlers.Store{}

	// subrouters
	userRouter := r.PathPrefix("/users").Subrouter()
	protectedRouter := r.PathPrefix("/protected").Subrouter()
	categoryRouter := r.PathPrefix("/categories").Subrouter()
	paymentRouter := r.PathPrefix("/payment").Subrouter()

	protectedRouter.Use(middlewares.JWTMiddleware)
	categoryRouter.Use(middlewares.JWTMiddleware)
	paymentRouter.Use(middlewares.JWTMiddleware)

	// users
	userRouter.HandleFunc("", store.GetUser).Methods(http.MethodGet)
	userRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		store.CreateUser(w, r)
	}).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/users/{id}", store.GetUser).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		store.UpdateUser(w, r)
	}).Methods(http.MethodPut)
	protectedRouter.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		store.DeleteUser(w, r)
	}).Methods(http.MethodDelete)

	// categories
	categoryRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		store.CreateCategory(w, r)
	}).Methods(http.MethodPost)
	categoryRouter.HandleFunc("", store.GetCategories).Methods(http.MethodGet)
	categoryRouter.HandleFunc("/{id}", store.GetCategories).Methods(http.MethodGet)
	categoryRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		store.UpdateCategory(w, r)
	}).Methods(http.MethodPut)
	categoryRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		store.DeleteCategory(w, r)
	}).Methods(http.MethodDelete)

	// payment
	paymentRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		store.CreatePayment(w, r)
	}).Methods(http.MethodPost)
	paymentRouter.HandleFunc("", store.GetPayment).Methods(http.MethodGet)
	paymentRouter.HandleFunc("/{id}", store.GetPayment).Methods(http.MethodGet)
	paymentRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		store.UpdatePayment(w, r)
	}).Methods(http.MethodPut)
	paymentRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		store.DeletePayment(w, r)
	}).Methods(http.MethodDelete)

	protectedRouter.HandleFunc("/email", func(w http.ResponseWriter, r *http.Request) {
		handlers.TestMail(w, r)
	}).Methods(http.MethodGet)

	return r
}
