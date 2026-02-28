package user

import (
	"context"
	"errors"

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

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (*entity.User, error) {
	if input.ID <= 0 {
		return nil, domainErr.ErrInvalidUserData
	}

	existingUser, err := uc.userRepo.GetUserByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrUserNotFound
		}

		return nil, err
	}

	profileChanged := false
	passwordChanged := false

	if input.Username != nil &&
		*input.Username != "" &&
		*input.Username != existingUser.Username {
		existingUser.Username = *input.Username
		profileChanged = true
	}

	if input.Email != nil &&
		*input.Email != "" &&
		*input.Email != existingUser.Email {
		if err := ValidateEmailFormat(*input.Email); err != nil {
			return nil, domainErr.ErrInvalidEmailFormat
		}

		existingUser.Email = *input.Email
		profileChanged = true
	}

	if input.Password != nil &&
		*input.Password != "" {
		hashedPassword, err := auth.HashPassword(*input.Password)
		if err != nil {
			return nil, domainErr.ErrPasswordHashing
		}

		existingUser.Password = hashedPassword
		passwordChanged = true
	}

	if profileChanged {
		if err := uc.userRepo.UpdateUserProfile(ctx, existingUser.ID, existingUser.Username, existingUser.Email); err != nil {
			return nil, err
		}
	}

	if passwordChanged {
		if err := uc.userRepo.UpdateUserPassword(ctx, existingUser.ID, existingUser.Password); err != nil {
			return nil, err
		}
	}

	return existingUser, nil
}
