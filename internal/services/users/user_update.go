package users

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	domainErrors "github.com/rober0xf/notifier/internal/domain/errors"
)

func (s Service) Update(user *domain.User) (*domain.User, error) {
	_, err := s.Repo.GetByID(user.ID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, err
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		return nil, dto.ErrInvalidUserData
	}

	hashed_password, err := authentication.HashPassword(user.Password)
	if err != nil {
		return nil, dto.ErrPasswordHashing
	}
	user.Password = hashed_password

	// update the user's fields using the repository
	if err := s.Repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
