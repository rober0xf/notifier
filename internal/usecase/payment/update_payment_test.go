package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/stretchr/testify/assert"

	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
)

func TestUpdatePayment(t *testing.T) {
	t.Run("successful update payment", func(t *testing.T) {
		uc, mockRepo := setupUpdatePaymentTest(t)

		payment := &entity.Payment{
			ID:        1,
			Name:      "Nike",
			Amount:    100.0,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeClothing,
			Date:      "2022-10-10",
			Paid:      false,
			Recurrent: false,
		}
		mockRepo.payments["1"] = payment

		newName := "Adidas"
		newAmount := 110.0
		input := UpdatePaymentInput{
			Name:   &newName,
			Amount: &newAmount,
		}

		updatedPayment, err := uc.Execute(context.Background(), 1, input)

		assert.NoError(t, err)
		assert.Equal(t, newName, updatedPayment.Name)
		assert.Equal(t, newAmount, updatedPayment.Amount)
		assert.Equal(t, payment.Type, updatedPayment.Type)
		assert.Equal(t, payment.Category, updatedPayment.Category)
		assert.Equal(t, payment.Date, updatedPayment.Date)
		assert.Equal(t, payment.Paid, updatedPayment.Paid)

		stored, _ := mockRepo.GetPaymentByID(context.Background(), 1)
		assert.Equal(t, newName, stored.Name)
		assert.Equal(t, newAmount, stored.Amount)
	})

	t.Run("successfully updates only name", func(t *testing.T) {
		uc, mockRepo := setupUpdatePaymentTest(t)

		payment := &entity.Payment{
			ID:       1,
			UserID:   1,
			Name:     "Original",
			Amount:   50.0,
			Type:     entity.TransactionTypeExpense,
			Category: entity.CategoryTypeElectronics,
			Date:     "2022-10-10",
		}
		mockRepo.payments["1"] = payment

		newName := "updated name"
		input := UpdatePaymentInput{
			Name: &newName,
		}

		updatedPayment, err := uc.Execute(context.Background(), 1, input)

		assert.NoError(t, err)
		assert.Equal(t, newName, updatedPayment.Name)
		assert.Equal(t, 50.0, updatedPayment.Amount)
	})

	t.Run("successfully updates recurrent payment frequency", func(t *testing.T) {
		uc, mockRepo := setupUpdatePaymentTest(t)

		monthly := entity.FrequencyTypeMonthly
		payment := &entity.Payment{
			ID:        1,
			UserID:    1,
			Name:      "NBA",
			Amount:    79.99,
			Type:      entity.TransactionTypeSubscription,
			Category:  entity.CategoryTypeEntertainment,
			Date:      "2022-10-10",
			Recurrent: true,
			Frequency: &monthly,
		}
		mockRepo.payments["1"] = payment

		yearly := entity.FrequencyTypeYearly
		input := UpdatePaymentInput{
			Frequency: &yearly,
		}

		updatedPayment, err := uc.Execute(context.Background(), 1, input)

		assert.NoError(t, err)
		assert.Equal(t, yearly, *updatedPayment.Frequency)
	})

	t.Run("returns error when payment not found", func(t *testing.T) {
		uc, _ := setupUpdatePaymentTest(t)

		name := "GPT"
		input := UpdatePaymentInput{
			Name: &name,
		}

		_, err := uc.Execute(context.Background(), 99999, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrPaymentNotFound)
	})

	t.Run("returns error for zero id", func(t *testing.T) {
		uc, _ := setupUpdatePaymentTest(t)

		name := "claude"
		input := UpdatePaymentInput{
			Name: &name,
		}

		_, err := uc.Execute(context.Background(), 0, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidPaymentData)
	})

	t.Run("returns error for negative id", func(t *testing.T) {
		uc, _ := setupUpdatePaymentTest(t)

		name := "claude"
		input := UpdatePaymentInput{
			Name: &name,
		}

		_, err := uc.Execute(context.Background(), -1, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidPaymentData)
	})

	t.Run("returns error for invalid amount", func(t *testing.T) {
		uc, mockRepo := setupUpdatePaymentTest(t)

		payment := &entity.Payment{
			ID:       1,
			UserID:   1,
			Name:     "claude",
			Type:     entity.TransactionTypeExpense,
			Category: entity.CategoryTypeEducation,
			Date:     "2026-03-21",
		}
		mockRepo.payments["1"] = payment

		invalidAmount := -120.0
		input := UpdatePaymentInput{
			Amount: &invalidAmount,
		}

		_, err := uc.Execute(context.Background(), -1, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidAmount)
	})

	t.Run("returns error for invalid transaction type", func(t *testing.T) {
		uc, mockRepo := setupUpdatePaymentTest(t)

		payment := &entity.Payment{
			ID:       1,
			UserID:   1,
			Name:     "codex",
			Amount:   100,
			Type:     entity.TransactionTypeExpense,
			Category: entity.CategoryTypeElectronics,
			Date:     "2022-10-10",
		}
		mockRepo.payments["1"] = payment

		invalidType := entity.TransactionType("invalid")
		input := UpdatePaymentInput{
			Type: &invalidType,
		}

		_, err := uc.Execute(context.Background(), 1, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidTransactionType)
	})

	t.Run("handles repository get errors", func(t *testing.T) {
		uc, mockRepo := setupUpdatePaymentTest(t)

		mockRepo.err = errors.New("database connection failed")

		name := "codex"
		input := UpdatePaymentInput{
			Name: &name,
		}

		_, err := uc.Execute(context.Background(), 1, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection failed")
	})

	t.Run("handles repository update errors", func(t *testing.T) {
		uc, mockRepo := setupUpdatePaymentTest(t)

		mockRepo.err = errors.New("database connection failed")

		payment := &entity.Payment{
			ID:       1,
			Name:     "leetcode",
			Amount:   10,
			Type:     entity.TransactionTypeSubscription,
			Category: entity.CategoryTypeEducation,
			Date:     "2022-11-11",
			Paid:     true,
		}
		mockRepo.payments["1"] = payment

		newName := "HBO"
		input := UpdatePaymentInput{
			Name: &newName,
		}

		mockRepo.err = errors.New("update failed")

		_, err := uc.Execute(context.Background(), 1, input)

		assert.Error(t, err)
	})
}
