package users

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/rober0xf/notifier/internal/domain/domain_errors"
)

func (s *Service) GetByEmail(email string) (*domain.User, error) {
	if email == "" {
		return nil, dto.ErrInvalidUserData
	}

	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, domain_errors.ErrNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s Service) GetAll() ([]*domain.User, error) {
	users, err := s.Repo.GetAllUsers()
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

func (s Service) GetByID(id uint) (*domain.User, error) {
	if id == 0 {
		return nil, dto.ErrInvalidUserData
	}

	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, domain_errors.ErrNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
