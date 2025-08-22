package payments

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Create(payment *domain.Payment) (*domain.Payment, error) {
	if payment == nil {
		return nil, errors.New("payment is nil")
	}

	if payment.UserID == 0 || payment.NetAmount <= 0 || payment.GrossAmount <= 0 || payment.Name == "" || payment.Type == "" || payment.Date.IsZero() {
		return nil, dto.ErrInvalidPaymentData
	}

	if err := s.Repo.CreatePayment(payment); err != nil {
		return nil, err
	}

	return payment, nil
}
