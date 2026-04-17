package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerification_Success_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	token := "verification-token"
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	user := &entity.User{
		Username:     "rober0xf",
		Email:        "rober0xf@gmail.com",
		PasswordHash: "hashedpassword",
		IsActive:     false,
	}
	_, err := deps.userRepo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	_, err = deps.userRepo.CreateUserToken(
		context.Background(),
		&entity.UserToken{
			UserID:    user.ID,
			TokenHash: tokenHash,
			Purpose:   entity.TokenPurposeEmailVerification,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/v1/users/email_verification/"+token, nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "email verified successfully", response["message"])

	updatedUser, err := deps.userRepo.GetUserByEmail(context.Background(), "rober0xf@gmail.com")
	require.NoError(t, err)
	assert.True(t, updatedUser.IsActive)
}

func TestVerification_NotFound_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	req := httptest.NewRequest("GET", "/v1/users/email_verification/abcd", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid or expired token", response["error"])
}

func TestVerification_Expired_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	token := "verification-token"
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	user := &entity.User{
		Username:     "rober0xf",
		Email:        "rober0xf@gmail.com",
		PasswordHash: "hashedpassword",
		IsActive:     false,
	}
	_, err := deps.userRepo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	// in the past
	_, err = deps.userRepo.CreateUserToken(
		context.Background(),
		&entity.UserToken{
			UserID:    user.ID,
			TokenHash: tokenHash,
			Purpose:   entity.TokenPurposeEmailVerification,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		})
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/v1/users/email_verification/"+token, nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid or expired verification link", response["error"])
}

func TestVerification_AlreadyVerified_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	token := "verification-token"
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	user := &entity.User{
		Username:     "rober0xf",
		Email:        "rober0xf@gmail.com",
		PasswordHash: "hashedpassword",
		IsActive:     false,
	}
	_, err := deps.userRepo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	_, err = deps.userRepo.CreateUserToken(
		context.Background(),
		&entity.UserToken{
			UserID:    user.ID,
			TokenHash: tokenHash,
			Purpose:   entity.TokenPurposeEmailVerification,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
	require.NoError(t, err)

	// activate before trying to verify again
	_, err = deps.userRepo.UpdateUserIsActiveReturning(context.Background(), user.ID, true)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/v1/users/email_verification/"+token, nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusConflict, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "email already verified", response["error"])
}
