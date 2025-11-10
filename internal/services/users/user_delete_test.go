package users

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/stretchr/testify/assert"
)

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int) error {
	if m.err != nil {
		return m.err
	}

	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return dto.ErrNotFound
}

func TestDeleteUser(t *testing.T) {
	t.Run("succesfully deleted user", func(t *testing.T) {
		service, mock_repo := setup_test_user_service(t)

		mock_repo.users["user@test.com"] = &domain.User{
			ID:       1,
			Email:    "user@test.com",
			Username: "testuser",
		}

		err := service.Delete(context.Background(), 1)
		assert.NoError(t, err)

		_, err = mock_repo.GetUserByID(context.Background(), 1)
		assert.True(t, errors.Is(err, dto.ErrNotFound))
	})
}
