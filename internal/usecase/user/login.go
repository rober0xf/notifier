package user

import (
	"context"
	"fmt"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/domain/repository"
	"github.com/rober0xf/notifier/pkg/auth"
)

type LoginUseCase struct {
	userRepo repository.UserRepository
	tokenGen auth.TokenGenerator
}

func NewLoginUseCase(
	userRepo repository.UserRepository,
	tokenGen auth.TokenGenerator,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo: userRepo,
		tokenGen: tokenGen,
	}
}

type LoginOutput struct {
	Token string
	User  *entity.User
}

func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (*LoginOutput, error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, auth.ErrInvalidCredentials
	}

	if !auth.VerifyPassword(user.PasswordHash, password) {
		return nil, auth.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, auth.ErrEmailNotVerified
	}

	token, err := uc.tokenGen.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("LoginUC.Execute failed to generate token: %w", err)
	}

	return &LoginOutput{Token: token, User: user}, nil
}
