package payment_test

import (
	"context"
	"strconv"

	"github.com/rober0xf/notifier/internal/domain/entity"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type MockPaymentRepository struct {
	payments map[string]*entity.Payment // key: id
	err      error
}

type MockUserRepositoryForPayment struct {
	users map[string]*entity.User // key: email or id
	err   error
}

func (m *MockPaymentRepository) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	if m.err != nil {
		return m.err
	}

	if payment.ID == 0 {
		payment.ID = int32(len(m.payments) + 1)
	}

	idStr := strconv.Itoa(int(payment.ID))
	if _, exists := m.payments[idStr]; exists {
		return repoErr.ErrAlreadyExists
	}

	m.payments[idStr] = payment

	return nil
}

func (m *MockPaymentRepository) GetPaymentByID(ctx context.Context, id int) (*entity.Payment, error) {
	if m.err != nil {
		return nil, m.err
	}

	idStr := strconv.Itoa(id)
	payment, exists := m.payments[idStr]
	if !exists {
		return nil, repoErr.ErrNotFound
	}

	paymentCopy := *payment
	return &paymentCopy, nil
}

func (m *MockPaymentRepository) GetAllPayments(ctx context.Context) ([]entity.Payment, error) {
	if m.err != nil {
		return nil, m.err
	}

	payments := make([]entity.Payment, 0, len(m.payments))
	for _, payment := range m.payments {
		payments = append(payments, *payment)
	}

	return payments, nil
}

func (m *MockPaymentRepository) GetAllPaymentsFromUser(ctx context.Context, userID int) ([]entity.Payment, error) {
	if m.err != nil {
		return nil, m.err
	}

	payments := make([]entity.Payment, 0)
	for _, payment := range m.payments {
		if payment.UserID == userID {
			payments = append(payments, *payment)
		}
	}

	if len(payments) == 0 {
		return []entity.Payment{}, nil
	}

	return payments, nil
}

func (m *MockPaymentRepository) UpdatePayment(ctx context.Context, payment *entity.Payment) error {
	if m.err != nil {
		return m.err
	}

	idStr := strconv.Itoa(int(payment.ID))
	_, exists := m.payments[idStr]
	if !exists {
		return repoErr.ErrNotFound
	}

	m.payments[idStr] = payment

	return nil
}

func (m *MockPaymentRepository) DeletePayment(ctx context.Context, id int) error {
	if m.err != nil {
		return m.err
	}

	idStr := strconv.Itoa(id)
	_, exists := m.payments[idStr]
	if !exists {
		return repoErr.ErrNotFound
	}

	delete(m.payments, idStr)

	return nil
}

func (m *MockUserRepositoryForPayment) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	user, exists := m.users[strconv.Itoa(userID)]
	if !exists {
		return nil, repoErr.ErrNotFound
	}
	userCopy := *user

	return &userCopy, nil
}
