package payment_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreatePayment(t *testing.T) {
	t.Run("successfully creates a new payment", func(t *testing.T) {
		uc, mockRepo := setupCreatePaymentTest(t)

		paidAt := "2026-03-12"
		input := &entity.Payment{
			UserID:    1,
			Name:      "copilot",
			Amount:    100,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeEducation,
			Date:      "2026-02-10",
			Paid:      true,
			Recurrent: false,
			PaidAt:    &paidAt,
		}

		payment, err := uc.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, "copilot", payment.Name)
		assert.Equal(t, 100.0, payment.Amount)
		assert.Equal(t, entity.TransactionTypeExpense, payment.Type)
		assert.Equal(t, entity.CategoryTypeEducation, payment.Category)
		assert.NotEqual(t, 0, payment.ID)

		storedPayment, err := mockRepo.GetPaymentByID(context.Background(), int(payment.ID))
		assert.NoError(t, err)
		assert.Equal(t, payment.Name, storedPayment.Name)
	})

	t.Run("returns error when payment already exists", func(t *testing.T) {
		uc, mockRepo := setupCreatePaymentTest(t)

		mockRepo.err = repoErr.ErrAlreadyExists
		payment := &entity.Payment{
			UserID:    1,
			Name:      "infra",
			Amount:    3500,
			Type:      entity.TransactionTypeIncome,
			Category:  entity.CategoryTypeWork,
			Date:      "2025-11-07",
			Paid:      false,
			Recurrent: false,
		}

		_, err := uc.Execute(context.Background(), payment)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrPaymentAlreadyExists)
	})

	t.Run("successfully creates recurrent payment with frequency", func(t *testing.T) {
		uc, _ := setupCreatePaymentTest(t)

		frequency := entity.FrequencyTypeWeekly
		paidAt := "2026-06-12"
		input := &entity.Payment{
			UserID:    1,
			Name:      "copilot",
			Amount:    100,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeEducation,
			Date:      "2026-02-10",
			Paid:      true,
			Recurrent: true,
			Frequency: &frequency,
			PaidAt:    &paidAt,
		}

		payment, err := uc.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.True(t, payment.Recurrent)
		assert.NotNil(t, payment.Frequency)
		assert.Equal(t, entity.FrequencyTypeWeekly, *payment.Frequency)
	})

	t.Run("returns error when missing userID", func(t *testing.T) {
		uc, _ := setupCreatePaymentTest(t)

		input := &entity.Payment{
			UserID:    0,
			Name:      "copilot",
			Amount:    100,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeEducation,
			Date:      "2026-02-10",
			Paid:      true,
			Recurrent: false,
		}

		payment, err := uc.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("returns error for empty name", func(t *testing.T) {
		uc, _ := setupCreatePaymentTest(t)

		input := &entity.Payment{
			UserID:    1,
			Name:      "",
			Amount:    100,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeEducation,
			Date:      "2026-02-10",
			Paid:      true,
			Recurrent: false,
		}

		payment, err := uc.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("returns error for invalid amount", func(t *testing.T) {
		uc, _ := setupCreatePaymentTest(t)

		input := &entity.Payment{
			UserID:    1,
			Name:      "copilot",
			Amount:    0,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeElectronics,
			Date:      "2026-02-10",
			Paid:      false,
			Recurrent: false,
		}

		payment, err := uc.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("returns error for negative amount", func(t *testing.T) {
		uc, _ := setupCreatePaymentTest(t)

		input := &entity.Payment{
			UserID:    1,
			Name:      "copilot",
			Amount:    -50,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeEntertainment,
			Date:      "2026-02-10",
			Paid:      true,
			Recurrent: false,
		}

		payment, err := uc.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("returns error for invalid transaction type", func(t *testing.T) {
		uc, _ := setupCreatePaymentTest(t)

		input := &entity.Payment{
			UserID:    1,
			Name:      "Test",
			Amount:    100,
			Type:      "invalid_type",
			Category:  entity.CategoryTypeClothing,
			Date:      "2026-02-10",
			Paid:      true,
			Recurrent: false,
		}

		payment, err := uc.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.ErrorIs(t, err, domainErr.ErrInvalidTransactionType)
	})

	t.Run("returns error for invalid category", func(t *testing.T) {
		uc, _ := setupCreatePaymentTest(t)

		input := &entity.Payment{
			UserID:    1,
			Name:      "Test",
			Amount:    100,
			Type:      entity.TransactionTypeExpense,
			Category:  "invalid_category",
			Date:      "2026-02-10",
			Paid:      true,
			Recurrent: false,
		}

		payment, err := uc.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.ErrorIs(t, err, domainErr.ErrInvalidCategory)
	})

	t.Run("returns error for missing date", func(t *testing.T) {
		uc, _ := setupCreatePaymentTest(t)

		input := &entity.Payment{
			UserID:   1,
			Name:     "Test",
			Amount:   100,
			Type:     entity.TransactionTypeExpense,
			Category: entity.CategoryTypeClothing,
			Date:     "",
		}

		payment, err := uc.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.Contains(t, err.Error(), "date is required")
	})

	t.Run("returns error for recurrent payment without frequency", func(t *testing.T) {
		uc, _ := setupCreatePaymentTest(t)

		input := &entity.Payment{
			UserID:    1,
			Name:      "Netflix",
			Amount:    15.99,
			Type:      entity.TransactionTypeSubscription,
			Category:  entity.CategoryTypeEntertainment,
			Date:      "2026-02-10",
			Recurrent: true,
			Frequency: nil,
		}

		payment, err := uc.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.ErrorIs(t, err, domainErr.ErrInvalidFrequency)
	})

	t.Run("handles repository errors", func(t *testing.T) {
		uc, mockRepo := setupCreatePaymentTest(t)

		mockRepo.err = errors.New("database connection failed")

		input := &entity.Payment{
			UserID:   1,
			Name:     "Test",
			Amount:   100,
			Type:     entity.TransactionTypeExpense,
			Category: entity.CategoryTypeClothing,
			Date:     "2026-02-10",
		}

		payment, err := uc.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, payment)
	})
}
