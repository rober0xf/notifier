package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/rober0xf/notifier/pkg/auth"
)

type UpdateUserUseCase struct {
	userRepo repository.UserRepository
}

func NewUpdateUserUseCase(userRepo repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo: userRepo,
	}
}

type UpdateUserInput struct {
	ID       int
	Username *string
	Email    *string
	Password *string
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput, userID int, role entity.Role) (*entity.User, error) {
	existingUser, err := uc.fetchUser(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if err := uc.authorizeUpdate(userID, existingUser.ID, role); err != nil {
		return nil, err
	}

	profileChanged, passwordChanged, err := uc.applyChanges(existingUser, input)
	if err != nil {
		return nil, err
	}

	if err := uc.persistChanges(ctx, existingUser, profileChanged, passwordChanged); err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (uc *UpdateUserUseCase) fetchUser(ctx context.Context, id int) (*entity.User, error) {
	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repoErr.ErrNotFound):
			return nil, domainErr.ErrUserNotFound
		default:
			return nil, fmt.Errorf("UpdateUserUC.fetchUser failed to get user by id: %w", err)
		}
	}

	return user, nil
}

func (uc *UpdateUserUseCase) authorizeUpdate(givenID, expectedID int, role entity.Role) error {
	if givenID != expectedID && role != entity.RoleAdmin {
		return auth.ErrForbidden
	}

	return nil
}

func (uc *UpdateUserUseCase) applyChanges(user *entity.User, input UpdateUserInput) (profileChanged, passwordChanged bool, err error) {
	if input.Username != nil && *input.Username != "" && *input.Username != user.Username {
		user.Username = *input.Username
		profileChanged = true
	}

	if input.Email != nil && *input.Email != "" && *input.Email != user.Email {
		if err := ValidateEmail(*input.Email, nil); err != nil {
			return false, false, domainErr.ErrInvalidEmailFormat
		}

		user.Email = *input.Email
		profileChanged = true
	}

	if input.Password != nil && *input.Password != "" {
		if err := ValidatePassword(*input.Password); err != nil {
			return false, false, domainErr.ErrInvalidPassword
		}

		hashedPassword, err := auth.HashPassword(*input.Password)
		if err != nil {
			return false, false, fmt.Errorf("UpdateUserUC.applyChanges failed to hash password: %w", err)
		}

		user.PasswordHash = hashedPassword
		passwordChanged = true
	}

	return profileChanged, passwordChanged, nil
}

func (uc *UpdateUserUseCase) persistChanges(ctx context.Context, user *entity.User, profileChanged, passwordChanged bool) error {
	if profileChanged {
		if err := uc.userRepo.UpdateUserProfile(ctx, user.ID, user.Username, user.Email); err != nil {
			switch {
			case errors.Is(err, repoErr.ErrNotFound):
				return domainErr.ErrUserNotFound
			case errors.Is(err, repoErr.ErrEmailAlreadyExists):
				return domainErr.ErrEmailAlreadyExists
			case errors.Is(err, repoErr.ErrUsernameAlreadyExists):
				return domainErr.ErrUsernameAlreadyExists
			default:
				return fmt.Errorf("UpdateUserUC.persistChanges failed to update user profile: %w", err)
			}
		}
	}

	if passwordChanged {
		if err := uc.userRepo.UpdateUserPassword(ctx, user.ID, user.PasswordHash); err != nil {
			switch {
			case errors.Is(err, repoErr.ErrNotFound):
				return domainErr.ErrUserNotFound
			default:
				return fmt.Errorf("UpdateUserUC.persistChanges failed to update user password: %w", err)
			}
		}
	}

	return nil
}
