package payment_test

import (
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/usecase/payment"
)

func setupTestPaymentRepo(t *testing.T) *MockPaymentRepository {
	t.Helper()

	return &MockPaymentRepository{
		payments: make(map[string]*entity.Payment),
	}
}

func setupTestUserRepo(t *testing.T) *MockUserRepositoryForPayment {
	t.Helper()

	return &MockUserRepositoryForPayment{
		users: make(map[string]*entity.User),
	}
}

func setupCreatePaymentTest(t *testing.T) (*payment.CreatePaymentUseCase, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)

	uc := payment.NewCreatePaymentUseCase(mockRepo)

	return uc, mockRepo
}

func setupGetPaymentByIDTest(t *testing.T) (*payment.GetPaymentByIDUseCase, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)
	uc := payment.NewGetPaymentByIDUseCase(mockRepo)

	return uc, mockRepo
}

func setupGetAllPaymentsTest(t *testing.T) (*payment.GetAllPaymentsUseCase, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)
	uc := payment.NewGetAllPaymentsUseCase(mockRepo)

	return uc, mockRepo
}

func setupGetAllPaymentsFromUserTest(t *testing.T) (*payment.GetAllPaymentsFromUserUseCase, *MockPaymentRepository, *MockUserRepositoryForPayment) {
	t.Helper()

	mockPaymentRepo := setupTestPaymentRepo(t)
	mockUserRepo := setupTestUserRepo(t)
	uc := payment.NewGetAllPaymentsFromUserUseCase(mockPaymentRepo, mockUserRepo)

	return uc, mockPaymentRepo, mockUserRepo
}

func setupDeletePaymentTest(t *testing.T) (*payment.DeletePaymentUseCase, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)
	uc := payment.NewDeletePaymentUseCase(mockRepo)

	return uc, mockRepo
}

func setupUpdatePaymentTest(t *testing.T) (*payment.UpdatePaymentUseCase, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)
	uc := payment.NewUpdatePaymentUseCase(mockRepo)

	return uc, mockRepo
}
