package users

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	domainErrors "github.com/rober0xf/notifier/internal/domain/errors"
)

func (s *Service) Create(name string, email string, password string) (*domain.User, error) {
	// check if the user already exists
	exists, err := s.Repo.GetByEmail(email)

	if err == nil && exists != nil {
		return nil, dto.ErrUserAlreadyExists
	} else if err != nil && !errors.Is(err, domainErrors.ErrNotFound) {
		return nil, err
	}

	hashed, err := authentication.HashPassword(password)
	if err != nil {
		return nil, dto.ErrPasswordHashing
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: hashed,
	}

	// store the user
	if err := s.Repo.Create(user); err != nil {
		if errors.Is(err, domainErrors.ErrAlreadyExists) {
			return nil, dto.ErrUserAlreadyExists
		}
		return nil, dto.ErrInternalServerError
	}

	return user, nil
}
