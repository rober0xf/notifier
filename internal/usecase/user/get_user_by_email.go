package user

import (
	"context"
	"errors"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type GetUserByEmailUseCase struct {
	userRepo repository.UserRepository
}

func NewGetUserByEmailUseCase(userRepo repository.UserRepository) *GetUserByEmailUseCase {
	return &GetUserByEmailUseCase{
		userRepo: userRepo,
	}
}

func (uc *GetUserByEmailUseCase) Execute(ctx context.Context, email string) (*entity.User, error) {
	if err := ValidateEmailFormat(email); err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}
