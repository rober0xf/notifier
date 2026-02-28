package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
	t.Run("succesfully deleted user", func(t *testing.T) {
		uc, mockRepo := setupDeleteUserTest(t)

		email := "richard@gmail.com"
		user := &entity.User{
			ID:       1,
			Username: "richard",
			Email:    email,
		}
		mockRepo.users[email] = user
		mockRepo.users["1"] = user

		err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)

		_, err = mockRepo.GetUserByID(context.Background(), 1)
		assert.ErrorIs(t, err, repoErr.ErrNotFound)

		_, err = mockRepo.GetUserByEmail(context.Background(), email)
		assert.ErrorIs(t, err, repoErr.ErrNotFound)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		uc, _ := setupDeleteUserTest(t)

		nonExistingID := 99999

		err := uc.Execute(context.Background(), nonExistingID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrUserNotFound)
	})

	t.Run("returns error for zero id", func(t *testing.T) {
		uc, mockRepo := setupDeleteUserTest(t)

		email := "richard@gmail.com"
		mockRepo.users[email] = &entity.User{
			ID:       0,
			Username: "richard",
			Email:    email,
		}

		err := uc.Execute(context.Background(), 0)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidUserData)
	})

	t.Run("returns error for negative id", func(t *testing.T) {
		uc, mockRepo := setupDeleteUserTest(t)

		email := "richard@gmail.com"
		mockRepo.users[email] = &entity.User{
			ID:       -1,
			Username: "richard",
			Email:    email,
		}

		err := uc.Execute(context.Background(), -1)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidUserData)
	})

	t.Run("handles repository error", func(t *testing.T) {
		uc, mockRepo := setupDeleteUserTest(t)

		mockRepo.err = errors.New("database connection failed")

		err := uc.Execute(context.Background(), 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection failed")
	})
}
