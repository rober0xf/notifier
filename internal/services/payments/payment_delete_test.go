package payments

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/stretchr/testify/assert"
)

func (m *MockPaymentRepository) DeletePayment(ctx context.Context, paymentID int) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.payments[paymentID]; !exists {
		return dto.ErrNotFound // return the repo error
	}
	delete(m.payments, paymentID)

	return nil
}

func TestDeletePayment(t *testing.T) {
	t.Run("successfully deleted payment", func(t *testing.T) {
		service, mock_repo := setup_test_payment_service(t)

		payment := &domain.Payment{
			ID:        1,
			UserID:    1,
			Name:      "Disney+",
			Amount:    100.0,
			Type:      domain.Subscription,
			Category:  domain.Entertainment,
			Date:      "2022-11-01",
			Paid:      true,
			Recurrent: true,
			Frequency: &monthly,
		}
		mock_repo.payments[1] = payment

		err := service.Delete(context.Background(), 1)
		assert.NoError(t, err)

		fetched_payment := mock_repo.payments[1]
		assert.Nil(t, fetched_payment)
	})
	t.Run("returns error when payment not found", func(t *testing.T) {
		service, _ := setup_test_payment_service(t)

		err := service.Delete(context.Background(), 99999)

		assert.Error(t, err)
		assert.Equal(t, err, dto.ErrPaymentNotFound)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		service, mock_repo := setup_test_payment_service(t)
		payment := &domain.Payment{
			ID:        1,
			UserID:    1,
			Name:      "Test",
			Amount:    50.0,
			Type:      domain.Expense,
			Category:  domain.Entertainment,
			Date:      "2024-07-01",
			Paid:      false,
			Recurrent: false,
		}
		mock_repo.payments[1] = payment
		expected_error := errors.New("internal server error")
		mock_repo.err = expected_error

		err := service.Delete(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, err, expected_error)
	})
}
