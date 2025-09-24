package users

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Update(user *domain.User) (*domain.User, error) {
	if user.ID <= 0 {
		return nil, dto.ErrInvalidUserData
	}

	existing, err := s.Repo.GetUserByID(user.ID)
	if err != nil {
		if errors.Is(err, dto.ErrRepository) {
			return nil, dto.ErrInternalServerError
		}
		if errors.Is(err, dto.ErrNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, err
	}

	// update only given fields, more like a patch request
	if user.Username != "" {
		existing.Username = user.Username
	}
	if user.Email != "" {
		if !validateEmail(user.Email) {
			return nil, dto.ErrInvalidUserData
		}
		existing.Email = user.Email
	}
	if user.Password != "" {
		hashed_password, err := authentication.HashPassword(user.Password)
		if err != nil {
			return nil, dto.ErrPasswordHashing
		}
		existing.Password = hashed_password
	}

	// update the user's fields using the repository
	if err := s.Repo.UpdateUser(existing); err != nil {
		return nil, err
	}

	return existing, nil
}
