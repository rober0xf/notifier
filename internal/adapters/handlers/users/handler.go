package users

import (
	"github.com/rober0xf/notifier/internal/services/auth"
	"github.com/rober0xf/notifier/internal/services/users"
)

type userHandler struct {
	UserService *users.Service
	Utils       *auth.Service
}

func NewUserHandler(service *users.Service, authService *auth.Service) *userHandler {
	return &userHandler{
		UserService: service,
		Utils:       authService,
	}
}
