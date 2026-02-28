package user

import (
	"context"
	"errors"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type GetUserByIDUseCase struct {
	userRepo repository.UserRepository
}

func NewGetUserByIDUseCase(userRepo repository.UserRepository) *GetUserByIDUseCase {
	return &GetUserByIDUseCase{
		userRepo: userRepo,
	}
}

func (uc *GetUserByIDUseCase) Execute(ctx context.Context, id int) (*entity.User, error) {
	if id <= 0 {
		return nil, domainErr.ErrInvalidUserData
	}

	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
