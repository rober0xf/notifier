package user

import (
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/pkg/email"
)

func setupTestRepo(t *testing.T) *MockUserRepository {
	t.Helper()

	return &MockUserRepository{
		users: make(map[string]*entity.User),
	}
}

func setupCreateUserTest(t *testing.T) (*CreateUserUseCase, *MockUserRepository, *email.MockSender) {
	t.Setenv("ENV", "test") // skip sending email
	t.Helper()

	mockRepo := setupTestRepo(t)

	mockEmailSender := email.NewMockSender()

	disposableDomains := []string{}
	baseURL := "http://localhost:3000"

	uc := NewCreateUserUseCase(mockRepo, mockEmailSender, disposableDomains, baseURL)

	return uc, mockRepo, mockEmailSender
}

func setupGetUserByEmailTest(t *testing.T) (*GetUserByEmailUseCase, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	uc := NewGetUserByEmailUseCase(mockRepo)

	return uc, mockRepo
}

func setupGetUserByIDTest(t *testing.T) (*GetUserByIDUseCase, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	uc := NewGetUserByIDUseCase(mockRepo)

	return uc, mockRepo
}

func setupGetAllUsersTest(t *testing.T) (*GetAllUsersUseCase, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	uc := NewGetAllUsersUseCase(mockRepo)

	return uc, mockRepo
}

func setupDeleteUserTest(t *testing.T) (*DeleteUserUseCase, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	uc := NewDeleteUserUseCase(mockRepo)

	return uc, mockRepo
}

func setupUpdateUserTest(t *testing.T) (*UpdateUserUseCase, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	uc := NewUpdateUserUseCase(mockRepo)

	return uc, mockRepo
}

func setupVerifyEmailTest(t *testing.T) (*VerifyEmailUseCase, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	uc := NewVerifyEmailUseCase(mockRepo)

	return uc, mockRepo
}
