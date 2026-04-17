package http

import (
	"github.com/rober0xf/notifier/internal/usecase/user"
	"github.com/rober0xf/notifier/pkg/auth"
)

type UserHandler struct {
	createUserUC     *user.CreateUserUseCase
	loginUC          *user.LoginUseCase
	getUserByIDUC    *user.GetUserByIDUseCase
	getUserByEmailUC *user.GetUserByEmailUseCase
	getAllUsersUC    *user.GetAllUsersUseCase
	updateUserUC     *user.UpdateUserUseCase
	deleteUserUC     *user.DeleteUserUseCase
	verifyEmailUC    *user.VerifyEmailUseCase
	oauthUC          *user.GoogleOAuthUseCase
	googleVerifier   auth.GoogleVerifier
}

func NewUserHandler(
	createUserUC *user.CreateUserUseCase,
	loginUC *user.LoginUseCase,
	getUserByIDUC *user.GetUserByIDUseCase,
	getUserByEmailUC *user.GetUserByEmailUseCase,
	getAllUsersUC *user.GetAllUsersUseCase,
	updateUserUC *user.UpdateUserUseCase,
	deleteUserUC *user.DeleteUserUseCase,
	verifyEmailUC *user.VerifyEmailUseCase,
	oauthUC *user.GoogleOAuthUseCase,
	googleVerifier auth.GoogleVerifier,
) *UserHandler {
	return &UserHandler{
		createUserUC:     createUserUC,
		loginUC:          loginUC,
		getUserByIDUC:    getUserByIDUC,
		getUserByEmailUC: getUserByEmailUC,
		getAllUsersUC:    getAllUsersUC,
		updateUserUC:     updateUserUC,
		deleteUserUC:     deleteUserUC,
		verifyEmailUC:    verifyEmailUC,
		oauthUC:          oauthUC,
		googleVerifier:   googleVerifier,
	}
}
