package users

import (
	"context"
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

func (s *Service) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if !validateEmail(email) {
		return nil, dto.ErrInvalidUserData
	}

	user, err := s.Repo.GetUserByEmail(ctx, email)
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

func (s Service) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.Repo.GetAllUsers(ctx)
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

	return users, nil
}

func (s Service) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := s.Repo.GetUserByID(ctx, id)
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
