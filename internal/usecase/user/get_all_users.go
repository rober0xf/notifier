package user

import (
	"context"

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
	return uc.userRepo.GetAllUsers(ctx)
}
