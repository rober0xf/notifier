package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/stretchr/testify/assert"

	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
)

func TestDeletePayment(t *testing.T) {
	t.Run("successfully deleted payment", func(t *testing.T) {
		uc, mockRepo := setupDeletePaymentTest(t)

		monthly := entity.FrequencyTypeMonthly
		payment := &entity.Payment{
			ID:        1,
			UserID:    1,
			Name:      "fight pass",
			Amount:    110,
			Type:      entity.TransactionTypeSubscription,
			Category:  entity.CategoryTypeEntertainment,
			Date:      "2022-11-01",
			Paid:      true,
			Recurrent: true,
			Frequency: &monthly,
		}
		mockRepo.payments["1"] = payment

		err := uc.Execute(context.Background(), 1)
		assert.NoError(t, err)

		_, exists := mockRepo.payments["1"]
		assert.Nil(t, exists)
	})

	t.Run("returns error when payment not found", func(t *testing.T) {
		uc, _ := setupDeletePaymentTest(t)

		nonExistingID := 99999
		err := uc.Execute(context.Background(), nonExistingID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrPaymentNotFound)
	})

	t.Run("returns error for id zero", func(t *testing.T) {
		uc, _ := setupDeletePaymentTest(t)

		err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidPaymentData)
	})

	t.Run("returns error for id negative", func(t *testing.T) {
		uc, _ := setupDeletePaymentTest(t)

		err := uc.Execute(context.Background(), -1)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidPaymentData)
	})

	t.Run("handles repository errors", func(t *testing.T) {
		uc, mockRepo := setupDeletePaymentTest(t)

		mockRepo.err = errors.New("database connection failed")

		payment := &entity.Payment{
			ID:        1,
			UserID:    1,
			Name:      "neetcode",
			Amount:    199,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeEducation,
			Date:      "2024-07-01",
			Paid:      false,
			Recurrent: false,
		}
		mockRepo.payments["1"] = payment

		err := uc.Execute(context.Background(), 1)

		assert.Error(t, err)
		assert.Contains(t, err, "database connection failed")

		_, exists := mockRepo.payments["1"]
		assert.True(t, exists)
	})
}
