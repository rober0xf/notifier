package payments

import (
	"errors"
	"time"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Update(id int, payment *domain.UpdatePayment) (*domain.Payment, error) {
	existing, err := s.Repo.GetPaymentByID(id)
	if err != nil {
		if errors.Is(err, dto.ErrRepository) {
			return nil, dto.ErrInternalServerError
		}
		if errors.Is(err, dto.ErrNotFound) {
			return nil, dto.ErrPaymentNotFound
		}
		return nil, err
	}

	/* check fields from the json */
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

	if payment.Recurrent != nil {
		existing.Recurrent = *payment.Recurrent
	}
	if payment.Frequency != nil {
		existing.Frequency = payment.Frequency
	}
	if payment.ReceiptURL != nil {
		existing.ReceiptURL = payment.ReceiptURL
	}
	/* --------------------------- */

	if err := s.Repo.UpdatePayment(existing); err != nil {
		switch {
		case errors.Is(err, dto.ErrPaymentNotFound):
			return nil, dto.ErrPaymentNotFound
		case errors.Is(err, dto.ErrRepository):
			return nil, dto.ErrInternalServerError
		default:
			return nil, dto.ErrInternalServerError
		}
	}
	return existing, nil
}
