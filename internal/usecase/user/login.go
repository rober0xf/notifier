package user

import (
	"context"

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

func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (*entity.User, error) {
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

	return user, nil
}
