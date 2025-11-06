package users

import (
	"context"
	"strconv"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/stretchr/testify/assert"
)

/* used by other users service too */
const jwt_for_test = "2f1f31a2a68b999d6cc40169bac7496a0def0ff6a33e3acb2f78e5487e480c41"

type MockUserRepository struct {
	// key: id or email
	users map[string]*domain.User
	err   error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func setup_test_user_service(t *testing.T) (*Service, *MockUserRepository) {
	t.Helper()
	mock_repo := NewMockUserRepository()
	service := NewUserService(mock_repo, []byte(jwt_for_test))
	return service, mock_repo
}

/* end */

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.users[user.Email]; exists {
		return dto.ErrAlreadyExists // return the repo error
	}
	if user.ID == 0 {
		user.ID = len(m.users) + 1
	}

	m.users[user.Email] = user
	m.users[strconv.Itoa(user.ID)] = user

	return nil
}

func TestCreateUser(t *testing.T) {
	t.Run("succesfully creates a new user", func(t *testing.T) {

		service, mock_repo := setup_test_user_service(t)

		email := "test@gmail.com"
		username := ""
		password := "password123#-!"

		user, err := service.Create(context.Background(), username, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, username, user.Username)
		assert.NotEqual(t, password, user.Password) // it should be hashed
		assert.NotEmpty(t, user.ID)

		stored_user, err := mock_repo.GetUserByID(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, stored_user)
		assert.Equal(t, user.ID, stored_user.ID)
		assert.Equal(t, user.Email, stored_user.Email)
	})
}
