package payment

import (
	"context"
	"errors"
	"time"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type UpdatePaymentUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewUpdatePaymentUseCase(paymentRepo repository.PaymentRepository) *UpdatePaymentUseCase {
	return &UpdatePaymentUseCase{
		paymentRepo: paymentRepo,
	}
}

type UpdatePaymentInput struct {
	Name       *string                 `json:"name,omitempty"`
	Amount     *float64                `json:"amount,omitempty"`
	Type       *entity.TransactionType `json:"type,omitempty"`
	Category   *entity.CategoryType    `json:"category,omitempty"`
	Date       *string                 `json:"date,omitempty"`
	DueDate    *string                 `json:"due_date,omitempty"`
	Paid       *bool                   `json:"paid"`
	PaidAt     *string                 `json:"paid_at,omitempty"`
	Recurrent  *bool                   `json:"recurrent,omitempty"`
	Frequency  *entity.FrequencyType   `json:"frequency,omitempty"`
	ReceiptURL *string                 `json:"receipt_url,omitempty"`
}

func (uc *UpdatePaymentUseCase) Execute(ctx context.Context, id int, input UpdatePaymentInput) (*entity.Payment, error) {
	if id <= 0 {
		return nil, domainErr.ErrInvalidPaymentData
	}

	existingPayment, err := uc.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrPaymentNotFound
		}
		return nil, err
	}

	// partial updates
	if err := applyPaymentUpdates(&input, existingPayment); err != nil {
		return nil, err
	}

	if err := ValidatePayment(existingPayment); err != nil {
		return nil, err
	}

	if err := uc.paymentRepo.UpdatePayment(ctx, existingPayment); err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrPaymentNotFound
		}
		return nil, err
	}

	return existingPayment, nil
}

func applyPaymentUpdates(input *UpdatePaymentInput, existing *entity.Payment) error {
	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Amount != nil {
		if *input.Amount <= 0 {
			return domainErr.ErrInvalidAmount
		}
		existing.Amount = *input.Amount
	}
	if input.Type != nil {
		if !input.Type.IsValid() {
			return domainErr.ErrInvalidTransactionType
		}
		existing.Type = *input.Type
	}
	if input.Category != nil {
		if !input.Category.IsValid() {
			return domainErr.ErrInvalidCategory
		}
		existing.Category = *input.Category
	}
	if input.Date != nil {
		existing.Date = *input.Date
	}
	if input.DueDate != nil {
		existing.DueDate = input.DueDate
	}
	if input.Recurrent != nil {
		existing.Recurrent = *input.Recurrent
	}
	if input.Frequency != nil {
		if !input.Frequency.IsValid() {
			return domainErr.ErrInvalidFrequency
		}
		existing.Frequency = input.Frequency
	}
	if input.ReceiptURL != nil {
		existing.ReceiptURL = input.ReceiptURL
	}

	if input.Paid != nil {
		existing.Paid = *input.Paid
		if *input.Paid {
			if existing.PaidAt == nil || *existing.PaidAt == "" {
				now := time.Now().Format("2006-01-02")
				existing.PaidAt = &now
			}
		} else {
			existing.PaidAt = nil
		}
	}

	if existing.Recurrent && existing.Frequency == nil {
		return domainErr.ErrInvalidFrequency
	}

	return nil
}
