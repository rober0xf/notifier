package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByEmail_Success_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	email := "rober0xf2@gmail.com"
	_ = createTestUser(t, deps, "rober0xf", email, "password1#!")

	token := getAuthToken(t, deps.router)

	req := httptest.NewRequest("GET", "/v1/auth/users/email/"+email, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	require.Equal(t, "rober0xf", response["username"])
	require.Equal(t, email, response["email"])
	assert.NotContains(t, response, "password")
}

func TestGetByEmail_Empty_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)
	token := getAuthToken(t, deps.router)

	req := httptest.NewRequest("GET", "/v1/auth/users/email", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "email parameter required")
}

func TestGetByEmail_Invalid_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)
	token := getAuthToken(t, deps.router)

	invalidEmail := "richard.com"
	req := httptest.NewRequest("GET", "/v1/auth/users/email/"+invalidEmail, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid email format")
}

func TestGetByEmail_NotFound_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)
	token := getAuthToken(t, deps.router)

	nonExistingEmail := "doesnotexists@example.com"
	req := httptest.NewRequest("GET", "/v1/auth/users/email/"+nonExistingEmail, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "user not found")
}
