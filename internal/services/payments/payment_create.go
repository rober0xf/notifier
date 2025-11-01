package payments

import (
	"context"
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	if err := s.Repo.CreatePayment(ctx, payment); err != nil {
		switch {
		case errors.Is(err, dto.ErrAlreadyExists):
			return nil, dto.ErrPaymentAlreadyExists
		default:
			return nil, dto.ErrInternalServerError
		}
	}

	return payment, nil
}
