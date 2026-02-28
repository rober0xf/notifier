package user_test

import (
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/usecase/user"
	"github.com/rober0xf/notifier/pkg/email"
)

func setupTestRepo(t *testing.T) *MockUserRepository {
	t.Helper()

	return &MockUserRepository{
		users: make(map[string]*entity.User),
	}
}

func setupCreateUserTest(t *testing.T) (*user.CreateUserUseCase, *MockUserRepository, *email.MockSender) {
	t.Helper()

	mockRepo := setupTestRepo(t)

	mockEmailSender := email.NewMockSender()

	disposableDomains := []string{}
	baseURL := "http://localhost:3000"

	uc := user.NewCreateUserUseCase(mockRepo, mockEmailSender, disposableDomains, baseURL)

	return uc, mockRepo, mockEmailSender
}

func setupGetUserByEmailTest(t *testing.T) (*user.GetUserByEmailUseCase, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	uc := user.NewGetUserByEmailUseCase(mockRepo)

	return uc, mockRepo
}

func setupGetUserByIDTest(t *testing.T) (*user.GetUserByIDUseCase, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	uc := user.NewGetUserByIDUseCase(mockRepo)

	return uc, mockRepo
}

func setupGetAllUsersTest(t *testing.T) (*user.GetAllUsersUseCase, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	uc := user.NewGetAllUsersUseCase(mockRepo)

	return uc, mockRepo
}

func setupDeleteUserTest(t *testing.T) (*user.DeleteUserUseCase, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	uc := user.NewDeleteUserUseCase(mockRepo)

	return uc, mockRepo
}

func setupUpdateUserTest(t *testing.T) (*user.UpdateUserUseCase, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	uc := user.NewUpdateUserUseCase(mockRepo)

	return uc, mockRepo
}

func setupVerifyEmailTest(t *testing.T) (*user.VerifyEmailUseCase, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	uc := user.NewVerifyEmailUseCase(mockRepo)

	return uc, mockRepo
}
