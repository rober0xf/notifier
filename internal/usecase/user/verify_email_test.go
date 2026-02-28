package user_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestVerifyEmail(t *testing.T) {
	t.Run("successfully activates when a verification is valid", func(t *testing.T) {
		uc, mockRepo := setupVerifyEmailTest(t)

		plainToken := "verification_123"
		tokenHash := sha256.Sum256([]byte(plainToken))
		tokenHashString := hex.EncodeToString(tokenHash[:])

		email := "richard@gmail.com"
		user := &entity.User{
			ID:                    1,
			Username:              "richard",
			Email:                 email,
			Password:              "hashedPassword",
			Active:                false,
			EmailVerificationHash: tokenHashString,
			TokenExpiresAt:        time.Now().Add(24 * time.Hour),
		}
		mockRepo.users[email] = user
		mockRepo.users["1"] = user

		verifiedUser, err := uc.Execute(context.Background(), email, plainToken)

		assert.NoError(t, err)
		assert.True(t, verifiedUser.Active)

		storedUser, err := mockRepo.GetUserByID(context.Background(), 1)

		assert.True(t, storedUser.Active)
	})

	t.Run("returns error when token is invalid", func(t *testing.T) {
		uc, mockRepo := setupVerifyEmailTest(t)

		realToken := "real_token"
		tokenHash := sha256.Sum256([]byte(realToken))

		email := "richard@gmail.com"
		user := &entity.User{
			ID:                    1,
			Username:              "richard",
			Email:                 email,
			Active:                false,
			EmailVerificationHash: hex.EncodeToString(tokenHash[:]),
			TokenExpiresAt:        time.Now().Add(24 * time.Hour),
		}
		mockRepo.users[email] = user
		mockRepo.users["1"] = user

		wrongToken := "wrong_token"
		_, err := uc.Execute(context.Background(), email, wrongToken)

		assert.Error(t, err)
		assert.ErrorIs(t, err, auth.ErrInvalidToken)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		uc, _ := setupVerifyEmailTest(t)

		_, err := uc.Execute(context.Background(), "notfound@gmail.com", "no_token")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrUserNotFound)
	})
}
