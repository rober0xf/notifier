package payments

import (
	"context"
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (s *Service) Delete(ctx context.Context, id int) error {
	if err := s.Repo.DeletePayment(ctx, id); err != nil {
		switch {
		case errors.Is(err, dto.ErrNotFound):
			return dto.ErrPaymentNotFound
		default:
			return dto.ErrInternalServerError
		}
	}

	return nil
}
