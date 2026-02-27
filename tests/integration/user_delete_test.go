package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteUser_Success_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	userID := createTestUser(t, deps, "user1", "usermail1@gmail.com", "password1!#")
	userIDStr := strconv.Itoa(userID)

	req := httptest.NewRequest("DELETE", "/v1/users/"+userIDStr, nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)

	_, err := deps.userRepo.GetUserByID(context.Background(), userID)
	assert.ErrorIs(t, err, domainErr.ErrUserNotFound)
}

func TestDeleteUser_InvalidID_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	req := httptest.NewRequest("DELETE", "/v1/users/aa", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid id")
}

func TestDeleteUser_NotFound_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	req := httptest.NewRequest("DELETE", "/v1/users/9999", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "user not found")
}
