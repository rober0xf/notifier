package users

import (
	"context"
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (s *Service) Delete(ctx context.Context, id int) error {
	err := s.Repo.DeleteUser(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrNotFound):
			return dto.ErrUserNotFound
		case errors.Is(err, dto.ErrRepository):
			return dto.ErrInternalServerError
		default:
			return dto.ErrInternalServerError
		}
	}

	return nil
}
