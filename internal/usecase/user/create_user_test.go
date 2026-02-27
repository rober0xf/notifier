package user_test

import (
	"context"
	"errors"
	"testing"

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

	t.Run("sends email when not in test mode", func(t *testing.T) {
		t.Setenv("ENV", "development")
		uc, _, mockEmailSender := setupCreateUserTest(t)

		user, err := uc.Execute(context.Background(), "john", "john@example.com", "password123#-!")

		assert.NoError(t, err)
		assert.Equal(t, 1, len(mockEmailSender.SentEmails))

		sentEmail := mockEmailSender.SentEmails[0]
		assert.Equal(t, []string{"john@example.com"}, sentEmail.To)
		assert.Equal(t, "Verify account", sentEmail.Subject)
		assert.Contains(t, sentEmail.Body, "verify")
		assert.Contains(t, sentEmail.Body, user.Email)
	})

	t.Run("returns error when email sending fails", func(t *testing.T) {
		t.Setenv("ENV", "development")
		uc, _, mockEmailSender := setupCreateUserTest(t)

		// simulate error
		mockEmailSender.Err = errors.New("SMTP connection failed")

		user, err := uc.Execute(context.Background(), "john", "john@example.com", "password123#-!")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "error sending verification email")
	})
}
