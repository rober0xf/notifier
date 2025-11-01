package users

import (
	"context"
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	if user.ID <= 0 {
		return nil, dto.ErrInvalidUserData
	}

	existing, err := s.Repo.GetUserByID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, dto.ErrRepository) {
			return nil, dto.ErrInternalServerError
		}
		if errors.Is(err, dto.ErrNotFound) {
			return nil, dto.ErrUserNotFound
		}

		return nil, err
	}

	var profile_changed = false
	var password_changed = false

	if user.Username != "" && user.Username != existing.Username {
		existing.Username = user.Username
		profile_changed = true
	}
	if user.Email != "" && user.Email != existing.Email {
		if !validateEmail(user.Email) {
			return nil, dto.ErrInvalidUserData
		}
		existing.Email = user.Email
		profile_changed = true
	}
	if user.Password != "" {
		hashed_password, err := authentication.HashPassword(user.Password)
		if err != nil {
			return nil, dto.ErrPasswordHashing
		}
		existing.Password = hashed_password
		password_changed = true
	}

	if profile_changed {
		if err := s.Repo.UpdateUserProfile(ctx, existing.ID, existing.Username, existing.Email); err != nil {
			return nil, err
		}
	}
	if password_changed {
		if err := s.Repo.UpdateUserPassword(ctx, existing.ID, existing.Password); err != nil {
			return nil, err
		}
	}

	return existing, nil
}
