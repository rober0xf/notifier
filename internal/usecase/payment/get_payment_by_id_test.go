package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetPaymentByID(t *testing.T) {
	t.Run("succesfully found the payment by id", func(t *testing.T) {
		uc, mockRepo := setupGetPaymentByIDTest(t)

		expectedPayment := &entity.Payment{
			ID:        1,
			UserID:    1,
			Name:      "copilot",
			Amount:    100,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeEducation,
			Date:      "2026-02-10",
			Paid:      true,
			Recurrent: false,
		}
		mockRepo.payments["1"] = expectedPayment

		payment, err := uc.Execute(context.Background(), int(expectedPayment.ID))

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, expectedPayment.ID, payment.ID)
		assert.Equal(t, expectedPayment.Name, payment.Name)
		assert.Equal(t, expectedPayment.Amount, payment.Amount)
		assert.Equal(t, entity.TransactionTypeExpense, payment.Type)
		assert.Equal(t, entity.CategoryTypeEducation, payment.Category)
	})

	t.Run("returns error when payment not found", func(t *testing.T) {
		uc, _ := setupGetPaymentByIDTest(t)

		nonExistingID := 99999
		payment, err := uc.Execute(context.Background(), nonExistingID)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.ErrorIs(t, err, domainErr.ErrPaymentNotFound)
	})

	t.Run("returns error for zero id", func(t *testing.T) {
		uc, _ := setupGetPaymentByIDTest(t)

		payment, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.ErrorIs(t, err, domainErr.ErrInvalidPaymentData)
	})

	t.Run("returns error for negative id", func(t *testing.T) {
		uc, _ := setupGetPaymentByIDTest(t)

		payment, err := uc.Execute(context.Background(), -1)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.ErrorIs(t, err, domainErr.ErrInvalidPaymentData)
	})

	t.Run("handles repository error", func(t *testing.T) {
		uc, mockRepo := setupGetPaymentByIDTest(t)

		mockRepo.err = errors.New("database connection failed")
		payment, err := uc.Execute(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.Contains(t, err.Error(), "database connection failed")
	})
}
