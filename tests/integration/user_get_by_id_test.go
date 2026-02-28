package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserByID_Success_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)
	token := getAuthToken(t, deps.router)

	userID := createTestUser(t, deps, "rober0xf", "rober0xf2@gmail.com", "password1#!")
	userIDStr := strconv.Itoa(userID)

	req := httptest.NewRequest("GET", "/v1/auth/users/"+userIDStr, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "rober0xf", response["username"])
	assert.Equal(t, "rober0xf2@gmail.com", response["email"])
}

func TestGetUserByID_InvalidID_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)
	token := getAuthToken(t, deps.router)

	req := httptest.NewRequest("GET", "/v1/auth/users/aa", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid user data")
}

func TestGetUserByID_NegativeID_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)
	token := getAuthToken(t, deps.router)

	req := httptest.NewRequest("GET", "/v1/auth/users/-1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "id must be positive")
}

func TestGetUserByID_NotFound_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)
	token := getAuthToken(t, deps.router)

	req := httptest.NewRequest("GET", "/v1/auth/users/9999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "user not found")
}
