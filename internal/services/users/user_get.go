package users

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	domainErrors "github.com/rober0xf/notifier/internal/domain/errors"
)

func (s *Service) Get(email string) (*domain.User, error) {
	if email == "" {
		return nil, dto.ErrInvalidUserData
	}

	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, domainErrors.ErrNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s Service) GetAllUsers() ([]*domain.User, error) {
	users, err := s.Repo.GetAll()
	if err != nil {
		return nil, err
	}

	// cast []domain.User to []*domain.User
	userPointers := make([]*domain.User, len(users))
	for i := range users {
		userPointers[i] = &users[i]
	}

	return userPointers, nil
}

func (s Service) GetUserFromID(id uint) (*domain.User, error) {
	user, err := s.Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, domainErrors.ErrNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
