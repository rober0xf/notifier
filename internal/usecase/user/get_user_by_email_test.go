package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail(t *testing.T) {
	t.Run("succesfully found the user by email", func(t *testing.T) {
		uc, mockRepo := setupGetUserByEmailTest(t)

		email := "richard@gmail.com"
		expectedUser := &entity.User{
			ID:       1,
			Email:    email,
			Username: "pied",
			Password: "password123!-#",
			Active:   true,
		}
		mockRepo.users[email] = expectedUser
		mockRepo.users["1"] = expectedUser

		user, err := uc.Execute(context.Background(), email)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "pied", user.Username)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		uc, _ := setupGetUserByEmailTest(t)

		noExists := "notfound@example.com"
		user, err := uc.Execute(context.Background(), noExists)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, domainErr.ErrUserNotFound)
	})

	t.Run("returns error for invalid email", func(t *testing.T) {
		uc, _ := setupGetUserByEmailTest(t)

		invalidEmail := "invalid.com"

		user, err := uc.Execute(context.Background(), invalidEmail)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("returns error for empty email", func(t *testing.T) {
		uc, _ := setupGetUserByEmailTest(t)

		user, err := uc.Execute(context.Background(), "")

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("handles repository errors", func(t *testing.T) {
		uc, mockRepo := setupGetUserByEmailTest(t)

		mockRepo.err = errors.New("database connection failed")
		user, err := uc.Execute(context.Background(), "doesntexists@email.com")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "database connection failed")
	})
}
