package payment

import (
	"fmt"
	"strings"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/domain/errors"
)

func ValidatePayment(payment *entity.Payment) error {
	if payment == nil {
		return errors.ErrInvalidPaymentData
	}
	if payment.UserID == 0 {
		return fmt.Errorf("user_id is required: %w", errors.ErrInvalidPaymentData)
	}
	if strings.TrimSpace(payment.Name) == "" {
		return fmt.Errorf("name cannot be empty: %w", errors.ErrInvalidPaymentData)
	}
	if payment.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0: %w", errors.ErrInvalidPaymentData)
	}
	if !payment.Type.IsValid() {
		return fmt.Errorf("invalid transaction type: %w", errors.ErrInvalidTransactionType)
	}
	if !payment.Category.IsValid() {
		return fmt.Errorf("invalid category type: %w", errors.ErrInvalidCategory)
	}
	if payment.Date == "" {
		return fmt.Errorf("date is required: %w", errors.ErrInvalidDate)
	}

	if payment.Recurrent {
		if payment.Frequency == nil {
			return fmt.Errorf("frequency is required for recurrent payments: %w", errors.ErrInvalidFrequency)
		}
		if !payment.Frequency.IsValid() {
			return fmt.Errorf("invalid frequency: %w", errors.ErrInvalidFrequency)
		}
	}

	if payment.Paid && payment.PaidAt == nil {
		return fmt.Errorf("paid_at is required when paid is true: %w", errors.ErrInvalidPaymentData)
	}

	return nil
}
