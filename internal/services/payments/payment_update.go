package payments

import (
	"context"
	"errors"
	"time"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Update(ctx context.Context, id int, payment *domain.UpdatePayment) (*domain.Payment, error) {
	existing, err := s.Repo.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			return nil, dto.ErrPaymentNotFound
		}
		return nil, dto.ErrInternalServerError
	}

	// partial updates
	apply_payment_updates(payment, existing)

	if err := s.Repo.UpdatePayment(ctx, existing); err != nil {
		switch {
		case errors.Is(err, dto.ErrNotFound):
			return nil, dto.ErrPaymentNotFound
		default:
			return nil, dto.ErrInternalServerError
		}
	}

	return existing, nil
}

func apply_payment_updates(payment *domain.UpdatePayment, existing *domain.Payment) {
	if payment.Name != nil {
		existing.Name = *payment.Name
	}
	if payment.Amount != nil {
		existing.Amount = *payment.Amount
	}
	if payment.Type != nil {
		existing.Type = *payment.Type
	}
	if payment.Category != nil {
		existing.Category = *payment.Category
	}
	if payment.Date != nil {
		existing.Date = *payment.Date
	}
	if payment.DueDate != nil {
		existing.DueDate = payment.DueDate
	}
	if payment.Recurrent != nil {
		existing.Recurrent = *payment.Recurrent
	}
	if payment.Frequency != nil {
		existing.Frequency = payment.Frequency
	}
	if payment.ReceiptURL != nil {
		existing.ReceiptURL = payment.ReceiptURL
	}

	if payment.Paid != nil {
		existing.Paid = *payment.Paid
		if *payment.Paid {
			if existing.PaidAt == nil || *existing.PaidAt == "" {
				date := time.Now().Format("2006-01-02")
				existing.PaidAt = &date
			}
		} else {
			existing.PaidAt = nil
		}
	}
}
