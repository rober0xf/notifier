package users

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/stretchr/testify/assert"
)

/* GET BY ID */

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, exists := m.users[strconv.Itoa(id)]
	if !exists {
		return nil, dto.ErrNotFound // we return the repo error, not the service error
	}

	user_copy := *user // for update user
	return &user_copy, nil
}

func TestGetUserByID(t *testing.T) {
	t.Run("succesfully found the user by id", func(t *testing.T) {
		service, _ := setup_test_user_service(t)

		// first we create the user
		email := "test@gmail.com"
		username := "username"
		password := "password123#_2"

		created_user, err := service.Create(context.Background(), username, email, password)
		assert.NoError(t, err)

		user, err := service.GetByID(context.Background(), created_user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, created_user.ID, user.ID)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, username, user.Username)
		assert.NotEqual(t, password, user.Password)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		service, _ := setup_test_user_service(t)

		// id that doesnt exists
		nonExistentID := 99999
		user, err := service.GetByID(context.Background(), nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, dto.ErrUserNotFound))
	})

	t.Run("handles repository errors", func(t *testing.T) {
		service, mock_repo := setup_test_user_service(t)

		// repo error
		mock_repo.err = errors.New("database connection failed")
		user, err := service.GetByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

/* END GET BY ID */

/* GET BY EMAIL */

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, exists := m.users[email]
	if !exists {
		return nil, dto.ErrNotFound
	}

	user_copy := *user
	return &user_copy, nil
}

func TestGetUserByEmail(t *testing.T) {
	t.Run("succesfully found the user by email", func(t *testing.T) {
		service, mock_repo := setup_test_user_service(t)

		email := "test@gmail.com"
		username := "testuser"
		password := "password123=#D"
		mock_repo.users[email] = &domain.User{
			ID:       1,
			Email:    email,
			Username: username,
			Password: password,
		}

		user, err := service.GetByEmail(context.Background(), email)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, 1, user.ID)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		service, _ := setup_test_user_service(t)

		non_exists_user := "notfound@example.com"
		user, err := service.GetByEmail(context.Background(), non_exists_user)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, dto.ErrUserNotFound))
	})

	t.Run("handles repository errors", func(t *testing.T) {
		service, mock_repo := setup_test_user_service(t)

		mock_repo.err = errors.New("database connection failed")
		user, err := service.GetByEmail(context.Background(), "doesntexists@email.com")

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

/* END GET BY EMAIL */

/* GET ALL */

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	seen := make(map[int]bool)
	users := make([]domain.User, 0)
	for _, user := range m.users {
		seen[user.ID] = true
		users = append(users, *user)
	}

	return users, nil
}

func TestGetAllUsers(t *testing.T) {
	t.Run("succesfully found all users", func(t *testing.T) {
		service, mock_repo := setup_test_user_service(t)

		mock_repo.users["user1@test.com"] = &domain.User{ID: 1, Email: "user1@test.com", Username: "user1"}
		mock_repo.users["user2@test.com"] = &domain.User{ID: 2, Email: "user2@test.com", Username: "user2"}
		mock_repo.users["user3@test.com"] = &domain.User{ID: 3, Email: "user3@test.com", Username: "user3"}

		users, err := service.GetAll(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 3)

		userIDs := make(map[int]bool)
		for _, u := range users {
			userIDs[u.ID] = true
		}

		assert.True(t, userIDs[1])
		assert.True(t, userIDs[2])
		assert.True(t, userIDs[3])
	})

	t.Run("returns empty list when no users exist", func(t *testing.T) {
		service, _ := setup_test_user_service(t)

		users, err := service.GetAll(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Empty(t, users) // empty, not nil
	})

	t.Run("handles repository errors", func(t *testing.T) {
		service, mock_repo := setup_test_user_service(t)

		mock_repo.err = errors.New("database connection failed")

		users, err := service.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, users)
	})
}

/* END GET ALL */
