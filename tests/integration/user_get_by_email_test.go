package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserByEmail_Success_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAdminToken(t, userDeps)
	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)
	email, err := extractEmailFromToken(token)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/v1/admin/users/email/"+email, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response dto.UserPayload
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, userID, response.ID)
	assert.Equal(t, email, response.Email)
}

func TestGetUserByEmail_InvalidFormat_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAdminToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/users/email/abb", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid email format", response.Error)
}

func TestGetUserByEmail_NotFound_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAdminToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/users/email/nonexistent@gmail.com", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user not found", response.Error)
}

func TestGetUserByEmail_Unauthorized_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)

	req := httptest.NewRequest("GET", "/v1/admin/users/email/test@gmail.com", nil)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "no token provided", response.Error)
}

func TestGetUserByEmail_Forbidden_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/users/email/test@gmail.com", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "forbidden", response.Error)
}
