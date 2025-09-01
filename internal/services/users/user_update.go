package users

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/rober0xf/notifier/internal/domain/domain_errors"
)

func (s Service) Update(user *domain.User) (*domain.User, error) {
	existing, err := s.Repo.GetUserByID(user.ID)
	if err != nil {
		if errors.Is(err, domain_errors.ErrNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, err
	}

	// update only given fields, more like a patch request
	if user.Name != "" {
		existing.Name = user.Name
	}
	if user.Email != "" {
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
