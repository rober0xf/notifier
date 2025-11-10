package payments

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	if err := validate_payment(payment); err != nil {
		return nil, err
	}
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

func validate_payment(payment *domain.Payment) error {
	if payment.UserID == 0 {
		return fmt.Errorf("user_id is required: %w", dto.ErrInvalidPaymentData)
	}
	if strings.TrimSpace(payment.Name) == "" {
		return fmt.Errorf("name cannot be empty: %w", dto.ErrInvalidPaymentData)
	}
	if payment.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0: %w", dto.ErrInvalidPaymentData)
	}
	if payment.Date == "" {
		return fmt.Errorf("date is required: %w", dto.ErrInvalidPaymentData)
	}
	if payment.Recurrent && payment.Frequency == nil {
		return fmt.Errorf("frequency is required for recurrent payments: %w", dto.ErrInvalidPaymentData)
	}

	return nil
}
