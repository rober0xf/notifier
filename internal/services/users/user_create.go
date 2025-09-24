package users

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Create(username string, email string, password string) (*domain.User, error) {
	// check if the user already exists
	exists_user, err := s.Repo.GetUserByEmail(email)
	if err != nil && !errors.Is(err, dto.ErrUserNotFound) {
		return nil, err
	}
	if exists_user != nil {
		return nil, dto.ErrUserAlreadyExists
	}

	hashed, err := authentication.HashPassword(password)
	if err != nil {
		return nil, dto.ErrPasswordHashing
	}

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: hashed,
	}

	// store the user
	if err := s.Repo.CreateUser(user); err != nil {
		if errors.Is(err, dto.ErrAlreadyExists) {
			return nil, dto.ErrUserAlreadyExists
		}
		if errors.Is(err, dto.ErrRepository) {
			return nil, dto.ErrInternalServerError
		}
		return nil, err
	}

	return user, nil
}
