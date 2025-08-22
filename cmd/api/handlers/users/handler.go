package users

import (
	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/ports"
)

type Handler struct {
	UserService ports.UserService
	AuthUtils   httphelpers.AuthHelper
}

// func NewUserHandler(userService users.UserService, authUtils httphelpers.AuthHelper) *Handler {
// 	return &Handler{
// 		UserService: userService,
// 		AuthUtils:   authUtils,
// 	}
// }
