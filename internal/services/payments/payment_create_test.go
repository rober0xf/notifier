package payments

import (
	"context"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/stretchr/testify/assert"
)

/* used for others payments services too */
type MockPaymentRepository struct {
	payments map[int]*domain.Payment
	users    map[string]int
	err      error
}

func NewMockPaymentRepository() *MockPaymentRepository {
	return &MockPaymentRepository{
		payments: make(map[int]*domain.Payment),
		users:    make(map[string]int),
	}
}

func setup_test_payment_service(t *testing.T) (*Service, *MockPaymentRepository) {
	t.Helper()
	mock_repo := NewMockPaymentRepository()
	service := NewPayments(mock_repo)
	return service, mock_repo
}

/* end */

func (m *MockPaymentRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	if m.err != nil {
		return m.err
	}
	// simulate auto increment
	if payment.ID == 0 {
		payment.ID = int32(len(m.payments) + 1)
	}
	if _, exists := m.payments[int(payment.ID)]; exists {
		return dto.ErrAlreadyExists
	}
	m.payments[int(payment.ID)] = payment

	return nil
}

func TestCreatePayment(t *testing.T) {
	t.Run("succesfully creates a new payment", func(t *testing.T) {
		service, mock_repo := setup_test_payment_service(t)

		name := "Electricity Bill"
		amount := 75.50
		freq := domain.Monthly
		userID := 1

		payment := &domain.Payment{
			UserID:    userID,
			Name:      name,
			Amount:    amount,
			Type:      domain.Expense,
			Category:  domain.Entertainment,
			Frequency: &freq,
			Date:      "2024-11-07",
			Paid:      false,
			Recurrent: true,
		}

		created_payment, err := service.Create(context.Background(), payment)
		assert.NoError(t, err)
		assert.NotNil(t, created_payment)
		assert.Equal(t, userID, created_payment.UserID)
		assert.Equal(t, name, created_payment.Name)
		assert.Equal(t, amount, created_payment.Amount)
		assert.NotZero(t, created_payment.ID)

		stored_payment, exists := mock_repo.payments[int(payment.ID)]
		assert.True(t, exists)
		assert.Equal(t, created_payment, stored_payment)
	})
	t.Run("returns error when payment already exists", func(t *testing.T) {
		service, _ := setup_test_payment_service(t)

		payment := &domain.Payment{
			ID:        1,
			UserID:    1,
			Name:      "Test Payment",
			Amount:    50.0,
			Type:      domain.Income,
			Category:  domain.Sports,
			Date:      "2025-11-07",
			Paid:      true,
			Recurrent: false,
		}
		_, err := service.Create(context.Background(), payment)
		assert.NoError(t, err)

		_, err = service.Create(context.Background(), payment)
		assert.Error(t, err)
		assert.Equal(t, dto.ErrPaymentAlreadyExists, err)
	})
}
