package user

import (
	"context"
	"fmt"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/domain/repository"
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
		return nil, fmt.Errorf("GetAllUsersUC.Execute failed to get all users: %w", err)
	}

	return users, nil
}
