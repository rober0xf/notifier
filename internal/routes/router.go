package routes

// func InitRouter(db *gorm.DB) *mux.Router {
// 	r := mux.NewRouter()
// 	store := &handlers.AuthHandler{}

// 	// subrouters
// 	_ = setupUserRoutes(r, store)
// 	_ = setupProtectedRoutes(r, store)
// 	_ = setupCategoryRoutes(r, store)
// 	_ = setupPaymentRoutes(r, store)

// 	return r
// }

// func setupUserRoutes(r *mux.Router, store *handlers.Store) *mux.Router {
// 	userRouter := r.PathPrefix("/api/users").Subrouter()

// 	userRouter.HandleFunc("", store.GetUser).Methods(http.MethodGet)
// 	userRouter.HandleFunc("", store.CreateUser).Methods(http.MethodPost)
// 	userRouter.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
// 		handlers.LoginHandler(w, r)
// 	}).Methods(http.MethodPost)

// 	return userRouter
// }

// func setupProtectedRoutes(r *mux.Router, store *handlers.Store) *mux.Router {
// 	protectedRouter := r.PathPrefix("api/auth").Subrouter()
// 	protectedRouter.Use(middlewares.JWTMiddleware)

// 	protectedRouter.HandleFunc("/users/{id}", store.GetUser).Methods(http.MethodGet)
// 	protectedRouter.HandleFunc("/users/{id}", store.UpdateUser).Methods(http.MethodPut)
// 	protectedRouter.HandleFunc("/users/{id}", store.DeleteUser).Methods(http.MethodDelete)

// 	protectedRouter.HandleFunc("/email", handlers.TestMail).Methods(http.MethodGet)

// 	return protectedRouter
// }

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
