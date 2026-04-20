package user_test

import (
	"context"
	"errors"
	"testing"

	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	t.Run("succesfully creates a new user", func(t *testing.T) {
		uc, _, mockEmailSender := setupCreateUserTest(t)

		user, err := uc.Execute(context.Background(), "richard", "piedpiper@gmail.com", "password123#-!")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, 0, len(mockEmailSender.SentEmails)) // 0 because env=test
		assert.NotEmpty(t, user.ID)
	})

	t.Run("fails when username in use", func(t *testing.T) {
		uc, _, _ := setupCreateUserTest(t)

		_, err := uc.Execute(context.Background(), "richard", "piedpiper@gmail.com", "password123#-!")
		assert.NoError(t, err)

		_, err = uc.Execute(context.Background(), "richard", "pied@gmail.com", "password123#-!")
		assert.Error(t, err)
		assert.Equal(t, err, repoErr.ErrUsernameAlreadyExists)
	})

	t.Run("fails when email in use", func(t *testing.T) {
		uc, _, _ := setupCreateUserTest(t)

		_, err := uc.Execute(context.Background(), "richard", "piedpiper@gmail.com", "password123#-!")
		assert.NoError(t, err)

		_, err = uc.Execute(context.Background(), "richard2", "piedpiper@gmail.com", "password123#-!")
		assert.Error(t, err)
		assert.Equal(t, err, repoErr.ErrEmailAlreadyExists)
	})

	t.Run("sends email when not in test mode", func(t *testing.T) {
		t.Setenv("ENV", "development")
		uc, _, mockEmailSender := setupCreateUserTest(t)

		mockEmailSender.ExpectSends(1)
		_, err := uc.Execute(context.Background(), "john", "john@example.com", "password123#-!")
		assert.NoError(t, err)

		mockEmailSender.Wait()
		assert.Equal(t, 1, len(mockEmailSender.SentEmails))

		sentEmail := mockEmailSender.SentEmails[0]
		assert.Equal(t, []string{"john@example.com"}, sentEmail.To)
		assert.Contains(t, sentEmail.Body, "verify account")
	})

	t.Run("creates user even when email sending fails", func(t *testing.T) {
		t.Setenv("ENV", "development")
		uc, _, mockEmailSender := setupCreateUserTest(t)

		// simulate error
		mockEmailSender.Err = errors.New("SMTP connection failed")
		mockEmailSender.ExpectSends(1)

		user, err := uc.Execute(context.Background(), "john", "john@example.com", "password123#-!")
		mockEmailSender.Wait()

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, 1, len(mockEmailSender.SentEmails))
	})
}
