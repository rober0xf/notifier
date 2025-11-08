package payments

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/stretchr/testify/assert"
)

var monthly = domain.Monthly

/* GET BY ID */
func (m *MockPaymentRepository) GetPaymentByID(ctx context.Context, paymentID int) (*domain.Payment, error) {
	if m.err != nil {
		return nil, m.err
	}
	payment, exists := m.payments[paymentID]
	if !exists {
		return nil, dto.ErrNotFound // return the repo error
	}

	return payment, nil
}

func TestGetPaymentByID(t *testing.T) {
	t.Run("succesfully found the payment by id", func(t *testing.T) {
		service, mock_repo := setup_test_payment_service(t)

		payment := &domain.Payment{
			UserID:    1,
			Name:      "Netflix",
			Amount:    25,
			Type:      domain.Subscription,
			Category:  domain.Entertainment,
			Date:      "2024-06-15",
			Paid:      false,
			Recurrent: false,
		}
		err := mock_repo.CreatePayment(context.Background(), payment)
		assert.NoError(t, err)

		found_payment, err := service.Get(context.Background(), int(payment.ID))

		assert.NoError(t, err)
		assert.NotNil(t, found_payment)
		assert.Equal(t, payment.ID, found_payment.ID)
		assert.Equal(t, "Netflix", found_payment.Name)
		assert.Equal(t, 25.00, found_payment.Amount)
		assert.Equal(t, 1, found_payment.UserID)
	})

	t.Run("returns error when payment not found", func(t *testing.T) {
		service, _ := setup_test_payment_service(t)

		nonExistentID := 99999
		payment, err := service.Get(context.Background(), nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.True(t, errors.Is(err, dto.ErrPaymentNotFound))
	})
}

/* END GET BY ID  */

/* GET ALL PAYMENTS */

func (m *MockPaymentRepository) GetAllPayments(ctx context.Context) ([]domain.Payment, error) {
	if m.err != nil {
		return nil, m.err
	}
	if len(m.payments) == 0 {
		return []domain.Payment{}, nil
	}
	payments := make([]domain.Payment, 0, len(m.payments))
	for _, payment := range m.payments {
		payments = append(payments, *payment)
	}

	return payments, nil
}

func TestGetAllPayments(t *testing.T) {
	t.Run("succesfully found all payments", func(t *testing.T) {
		service, mock_repo := setup_test_payment_service(t)

		payment1 := &domain.Payment{
			UserID:    1,
			Name:      "Cursor",
			Amount:    30.00,
			Type:      domain.Subscription,
			Category:  domain.Work,
			Date:      "2025-06-10",
			Paid:      true,
			Recurrent: true,
			Frequency: &monthly,
		}
		payment2 := &domain.Payment{
			UserID:    2,
			Name:      "Internship",
			Amount:    2000.00,
			Type:      domain.Income,
			Category:  domain.Work,
			Date:      "2025-06-01",
			Paid:      true,
			Recurrent: true,
			Frequency: &monthly,
		}
		err := mock_repo.CreatePayment(context.Background(), payment1)
		assert.NoError(t, err)
		err = mock_repo.CreatePayment(context.Background(), payment2)
		assert.NoError(t, err)

		payments, err := service.GetAllPayments(context.Background())

		assert.NoError(t, err)
		assert.Len(t, payments, 2)

		paymentMap := make(map[int32]domain.Payment)
		for _, p := range payments {
			paymentMap[p.ID] = p
		}

		p1, exists := paymentMap[payment1.ID]
		assert.True(t, exists)
		assert.Equal(t, "Cursor", p1.Name)
		assert.Equal(t, 30.00, p1.Amount)
		assert.Equal(t, domain.Subscription, p1.Type)

		p2, exists := paymentMap[payment2.ID]
		assert.True(t, exists)
		assert.Equal(t, "Internship", p2.Name)
		assert.Equal(t, 2000.00, p2.Amount)
		assert.Equal(t, domain.Income, p2.Type)
	})

	t.Run("returns empty list when no payments exist", func(t *testing.T) {
		service, _ := setup_test_payment_service(t)

		payments, err := service.GetAllPayments(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, payments)
		assert.Empty(t, payments)
		assert.Len(t, payments, 0)
	})
}

/* END GET ALL PAYMENTS */

/* GET ALL PAYMENTS FROM USER */
func (m *MockPaymentRepository) GetAllPaymentsFromUser(ctx context.Context, email string) ([]domain.Payment, error) {
	if m.err != nil {
		return nil, m.err
	}

	userID, ok := m.users[email]
	if !ok {
		return []domain.Payment{}, nil
	}
	payments := make([]domain.Payment, 0)
	for _, p := range m.payments {
		if p.UserID == userID {
			payments = append(payments, *p)
		}
	}

	return payments, nil
}

func (m *MockPaymentRepository) AddUser(email string, userID int) {
	m.users[email] = userID
}
func TestGetAllPaymentsFromUser(t *testing.T) {
	t.Run("succesfully found all payments from user", func(t *testing.T) {
		service, mock_repo := setup_test_payment_service(t)

		mock_repo.AddUser("user1@test.com", 1)
		mock_repo.AddUser("user2@test.com", 2)
		payment1 := &domain.Payment{
			UserID:    1,
			Name:      "Spotify",
			Amount:    10.00,
			Type:      domain.Subscription,
			Category:  domain.Entertainment,
			Date:      "2024-07-01",
			Paid:      true,
			Recurrent: true,
			Frequency: &monthly,
		}
		payment2 := &domain.Payment{
			UserID:    1,
			Name:      "Freelance Project",
			Amount:    500.00,
			Type:      domain.Income,
			Category:  domain.Work,
			Date:      "2024-07-05",
			Paid:      true,
			Recurrent: false,
		}
		payment3 := &domain.Payment{
			UserID:    2,
			Name:      "Gym Membership",
			Amount:    40.00,
			Type:      domain.Subscription,
			Category:  domain.Sports,
			Date:      "2024-07-03",
			Paid:      true,
			Recurrent: true,
			Frequency: &monthly,
		}

		err := mock_repo.CreatePayment(context.Background(), payment1)
		assert.NoError(t, err)
		err = mock_repo.CreatePayment(context.Background(), payment2)
		assert.NoError(t, err)
		err = mock_repo.CreatePayment(context.Background(), payment3)
		assert.NoError(t, err)

		// user 1
		payments_1, err := service.GetAllPaymentsFromUser(context.Background(), "user1@test.com")
		assert.NoError(t, err)
		assert.Len(t, payments_1, 2)
		for _, p := range payments_1 {
			assert.Equal(t, 1, p.UserID)
		}
		paymentNames := make([]string, len(payments_1))
		for i, p := range payments_1 {
			paymentNames[i] = p.Name
		}
		assert.Contains(t, paymentNames, "Spotify")
		assert.Contains(t, paymentNames, "Freelance Project")

		// user 2
		payments_2, err := service.GetAllPaymentsFromUser(context.Background(), "user2@test.com")
		assert.NoError(t, err)
		assert.Len(t, payments_2, 1)
		assert.Equal(t, 2, payments_2[0].UserID)
		assert.Equal(t, "Gym Membership", payments_2[0].Name)
	})

	t.Run("returns empty list when user has no payments", func(t *testing.T) {
		service, _ := setup_test_payment_service(t)

		payments, err := service.GetAllPaymentsFromUser(context.Background(), "nouser@test.com")

		assert.NoError(t, err)
		assert.NotNil(t, payments)
		assert.Empty(t, payments)
		assert.Len(t, payments, 0)
	})
}

/* END GET ALL FROM USER */
