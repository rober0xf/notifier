package api

import (
	"os"

	routes "github.com/rober0xf/notifier/internal/delivery/http"
	"github.com/rober0xf/notifier/internal/domain/repository"
	"github.com/rober0xf/notifier/internal/usecase/user"
	"github.com/rober0xf/notifier/pkg/auth"
	"github.com/rober0xf/notifier/pkg/email"
)

func buildUserRoutes(userRepo repository.UserRepository, tokenGen auth.TokenGenerator, emailSender email.EmailSender, googleClientID string) *routes.UserHandler {
	return routes.NewUserHandler(
		user.NewCreateUserUseCase(userRepo, emailSender, email.MustDisposableEmail(), os.Getenv("BASE_URL")),
		user.NewLoginUseCase(userRepo, tokenGen),
		user.NewGetUserByIDUseCase(userRepo),
		user.NewGetUserByEmailUseCase(userRepo),
		user.NewGetAllUsersUseCase(userRepo),
		user.NewUpdateUserUseCase(userRepo),
		user.NewDeleteUserUseCase(userRepo),
		user.NewVerifyEmailUseCase(userRepo),
		user.NewOAuthUseCase(userRepo, tokenGen),
		auth.NewGoogleVerifier(googleClientID),
	)
}
