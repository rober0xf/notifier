package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rober0xf/notifier/internal/handlers"
	Mail "github.com/rober0xf/notifier/internal/handlers"
	"github.com/rober0xf/notifier/internal/middlewares"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()
	handlers := &handlers.AuthHandler{}

	// subrouters
	_ = setup_user_routes(r, handlers)
	_ = setup_protected_routes(r, handlers)

	// TODO: implement category and payment routes when Store type is defined
	// _ = setupCategoryRoutes(r, handlers)
	// _ = setupPaymentRoutes(r, handlers)

	return r
}

func setup_user_routes(r *mux.Router, handlers *handlers.AuthHandler) *mux.Router {
	user_routes := r.PathPrefix("/api/users").Subrouter()

	user_routes.HandleFunc("", handlers.GetAllUsersHandler).Methods(http.MethodGet)
	user_routes.HandleFunc("", handlers.CreateUserHandler).Methods(http.MethodPost)
	user_routes.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r)
	}).Methods(http.MethodPost)
	user_routes.HandleFunc("/mail", func(w http.ResponseWriter, r *http.Request) {
		Mail.TestMail(w, r)
	}).Methods(http.MethodGet)

	return user_routes
}

func setup_protected_routes(r *mux.Router, handlers *handlers.AuthHandler) *mux.Router {
	protected_routes := r.PathPrefix("api/auth").Subrouter()
	protected_routes.Use(middlewares.JWTMiddleware)

	protected_routes.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		handlers.GetUserByIDHandler(w, r, id)
	}).Methods(http.MethodGet)
	protected_routes.HandleFunc("/users/{id}", handlers.UpdateUserHandler).Methods(http.MethodPut)
	protected_routes.HandleFunc("/users/{id}", handlers.DeleteUserHandler).Methods(http.MethodDelete)

	protected_routes.HandleFunc("/email", func(w http.ResponseWriter, r *http.Request) {
		Mail.TestMail(w, r)
	}).Methods(http.MethodGet)

	return protected_routes
}

// TODO:
// func setupCategoryRoutes(r *mux.Router, store *handlers.Store) *mux.Router {
// 	categoryRouter := r.PathPrefix("/api/categories").Subrouter()
// 	categoryRouter.Use(middlewares.JWTMiddleware)

// 	categoryRouter.HandleFunc("", store.GetCategories).Methods(http.MethodGet)
// 	categoryRouter.HandleFunc("", store.CreateCategory).Methods(http.MethodPost)
// 	categoryRouter.HandleFunc("/{id}", store.GetCategories).Methods(http.MethodGet)
// 	categoryRouter.HandleFunc("/{id}", store.UpdateCategory).Methods(http.MethodPut)
// 	categoryRouter.HandleFunc("/{id}", store.DeleteCategory).Methods(http.MethodDelete)

// 	return categoryRouter
// }

// func setupPaymentRoutes(r *mux.Router, store *handlers.Store) *mux.Router {
// 	paymentRouter := r.PathPrefix("/api/payments").Subrouter()
// 	paymentRouter.Use(middlewares.JWTMiddleware)

// 	paymentRouter.HandleFunc("", store.GetPayment).Methods(http.MethodGet)
// 	paymentRouter.HandleFunc("", store.CreatePayment).Methods(http.MethodPost)
// 	paymentRouter.HandleFunc("/{id}", store.GetPayment).Methods(http.MethodGet)
// 	paymentRouter.HandleFunc("/{id}", store.UpdatePayment).Methods(http.MethodPut)
// 	paymentRouter.HandleFunc("/{id}", store.DeletePayment).Methods(http.MethodDelete)

// 	return paymentRouter
// }
