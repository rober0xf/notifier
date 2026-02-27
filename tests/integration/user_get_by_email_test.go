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
	deps := setupTestUserDependencies(t)

	email := "rober0xf@gmail.com"
	_ = createTestUser(t, deps, "rober0xf", email, "password1#!")

	req := httptest.NewRequest("GET", "/v1/users/email/"+email, nil)
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
	deps := setupTestUserDependencies(t)

	req := httptest.NewRequest("GET", "/v1/users/email", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "email parameter required")
}

func TestGetByEmail_Invalid_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	invalidEmail := "richard.com"
	req := httptest.NewRequest("GET", "/v1/users/email/"+invalidEmail, nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid email format")
}

func TestGetByEmail_NotFound_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	nonExistingEmail := "doesnotexists@example.com"
	req := httptest.NewRequest("GET", "/v1/users/email/"+nonExistingEmail, nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "user not found")
}
