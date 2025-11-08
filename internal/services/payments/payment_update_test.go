package payments

import (
	"context"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/stretchr/testify/assert"
)

func (m *MockPaymentRepository) UpdatePayment(ctx context.Context, payment *domain.Payment) error {
	if m.err != nil {
		return m.err
	}
	if payment.ID == 0 {
		return dto.ErrInvalidData
	}
	if _, exists := m.payments[int(payment.ID)]; !exists {
		return dto.ErrNotFound
	}
	m.payments[int(payment.ID)] = payment

	return nil
}

func TestUpdatePayment(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		service, mock_repo := setup_test_payment_service(t)

		payment_id := 1
		payment := &domain.Payment{
			ID:       int32(payment_id),
			Name:     "Copilot",
			Amount:   100.0,
			Type:     domain.Subscription,
			Category: domain.Work,
			Date:     "2022-10-10",
			Paid:     false,
		}
		mock_repo.payments[payment_id] = payment

		new_name := "Copilot+"
		new_amount := 150.0

		updated_payment := &domain.UpdatePayment{
			Name:   &new_name,
			Amount: &new_amount,
		}
		new_payment, err := service.Update(context.Background(), 1, updated_payment)
		assert.NoError(t, err)
		assert.Equal(t, new_name, new_payment.Name)
		assert.Equal(t, new_amount, new_payment.Amount)
		assert.Equal(t, payment.Type, new_payment.Type)
		assert.Equal(t, payment.Category, new_payment.Category)
		assert.Equal(t, payment.Date, new_payment.Date)
		assert.Equal(t, payment.Paid, new_payment.Paid)
	})

	t.Run("update non-existent payment", func(t *testing.T) {
		service, _ := setup_test_payment_service(t)
		new_name := "Copilot+"
		updated_payment := &domain.UpdatePayment{
			Name: &new_name,
		}

		_, err := service.Update(context.Background(), 999, updated_payment)
		assert.Error(t, err)
		assert.Equal(t, dto.ErrPaymentNotFound, err)
	})

	t.Run("repository error on update", func(t *testing.T) {
		service, mock_repo := setup_test_payment_service(t)
		payment_id := 1
		payment := &domain.Payment{
			ID:       int32(payment_id),
			Name:     "Netflix",
			Amount:   15.0,
			Type:     domain.Subscription,
			Category: domain.Entertainment,
			Date:     "2022-11-11",
			Paid:     true,
		}
		mock_repo.payments[payment_id] = payment
		mock_repo.err = dto.ErrInternalServerError

		new_name := "HBO"
		updated_payment := &domain.UpdatePayment{
			Name: &new_name,
		}

		_, err := service.Update(context.Background(), 1, updated_payment)
		assert.Error(t, err)
		assert.Equal(t, err, dto.ErrInternalServerError)
		assert.Equal(t, "Netflix", mock_repo.payments[payment_id].Name)
	})
}
