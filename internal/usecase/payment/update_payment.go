package payment

import (
	"context"
	"errors"
	"fmt"
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
	Name       *string
	Amount     *float64
	Type       *entity.TransactionType
	Category   *entity.CategoryType
	Date       *string
	DueDate    *string
	Paid       *bool
	PaidAt     *string
	Recurrent  *bool
	Frequency  *entity.FrequencyType
	ReceiptURL *string
}

func (uc *UpdatePaymentUseCase) Execute(ctx context.Context, id int, input UpdatePaymentInput) error {
	if id <= 0 {
		return domainErr.ErrInvalidPaymentData
	}

	existingPayment, err := uc.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return domainErr.ErrPaymentNotFound
		}

		return fmt.Errorf("failed to get payment by id: %w", err)
	}

	// partial updates
	applyPaymentUpdates(&input, existingPayment)

	if err := uc.paymentRepo.UpdatePayment(ctx, existingPayment); err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return domainErr.ErrPaymentNotFound
		}

		return fmt.Errorf("UpdatePaymentUC.Execute failed to update payment: %w", err)
	}

	return nil
}

func applyPaymentUpdates(input *UpdatePaymentInput, existing *entity.Payment) {
	if input.Name != nil {
		existing.Name = *input.Name
	}

	if input.Amount != nil {
		existing.Amount = *input.Amount
	}

	if input.Type != nil {
		existing.Type = *input.Type
	}

	if input.Category != nil {
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
}
