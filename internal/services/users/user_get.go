package users

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func validateEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (s *Service) GetByEmail(email string) (*domain.User, error) {
	if !validateEmail(email) {
		return nil, dto.ErrInvalidUserData
	}
	return s.Repo.GetUserByEmail(email)
}

func (s Service) GetAll() ([]domain.User, error) {
	return s.Repo.GetAllUsers()
}

func (s Service) GetByID(id int) (*domain.User, error) {
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrNotFound):
			return nil, dto.ErrUserNotFound
		case errors.Is(err, dto.ErrRepository):
			return nil, dto.ErrInternalServerError
		default:
			return nil, dto.ErrInternalServerError
		}
	}
	return user, nil
}
