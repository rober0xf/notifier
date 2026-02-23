package payment

import (
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
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

func setupCreatePaymentTest(t *testing.T) (*CreatePaymentUseCase, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)

	uc := NewCreatePaymentUseCase(mockRepo)

	return uc, mockRepo
}

func setupGetPaymentByIDTest(t *testing.T) (*GetPaymentByIDUseCase, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)
	uc := NewGetPaymentByIDUseCase(mockRepo)

	return uc, mockRepo
}

func setupGetAllPaymentsTest(t *testing.T) (*GetAllPaymentsUseCase, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)
	uc := NewGetAllPaymentsUseCase(mockRepo)

	return uc, mockRepo
}

func setupGetAllPaymentsFromUserTest(t *testing.T) (*GetAllPaymentsFromUserUseCase, *MockPaymentRepository, *MockUserRepositoryForPayment) {
	t.Helper()

	mockPaymentRepo := setupTestPaymentRepo(t)
	mockUserRepo := setupTestUserRepo(t)
	uc := NewGetAllPaymentsFromUserUseCase(mockPaymentRepo, mockUserRepo)

	return uc, mockPaymentRepo, mockUserRepo
}

func setupDeletePaymentTest(t *testing.T) (*DeletePaymentUseCase, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)
	uc := NewDeletePaymentUseCase(mockRepo)

	return uc, mockRepo
}

func setupUpdatePaymentTest(t *testing.T) (*UpdatePaymentUseCase, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)
	uc := NewUpdatePaymentUseCase(mockRepo)

	return uc, mockRepo
}
