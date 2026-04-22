package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	authErr "github.com/rober0xf/notifier/pkg/auth"
)

type DeleteUserUseCase struct {
	userRepo repository.UserRepository
}

func NewDeleteUserUseCase(userRepo repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id, userID int, role entity.Role) error {
	if id <= 0 {
		return domainErr.ErrInvalidUserData
	}

	if id != userID && role != entity.RoleAdmin {
		return authErr.ErrForbidden
	}

	err := uc.userRepo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return domainErr.ErrUserNotFound
		}

		return fmt.Errorf("DeleteUserUC.Execute failed to delete user: %w", err)
	}

	return nil
}
