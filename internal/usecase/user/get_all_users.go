package user

import (
	"context"
	"errors"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type GetAllUsersUseCase struct {
	userRepo repository.UserRepository
}

func NewGetAllUsersUseCase(userRepo repository.UserRepository) *GetAllUsersUseCase {
	return &GetAllUsersUseCase{
		userRepo: userRepo,
	}
}

func (uc *GetAllUsersUseCase) Execute(ctx context.Context) ([]entity.User, error) {
	users, err := uc.userRepo.GetAllUsers(ctx)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return []entity.User{}, nil // return empty list
		}

		return nil, err
	}

	return users, nil
}
