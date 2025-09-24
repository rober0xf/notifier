package payments

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Create(payment *domain.Payment) (*domain.Payment, error) {
	if err := s.Repo.CreatePayment(payment); err != nil {
		switch {
		case errors.Is(err, dto.ErrPaymentAlreadyExists):
			return nil, dto.ErrPaymentAlreadyExists
		case errors.Is(err, dto.ErrRepository):
			return nil, dto.ErrInternalServerError
		default:
			return nil, dto.ErrInternalServerError
		}
	}
	return payment, nil
}
