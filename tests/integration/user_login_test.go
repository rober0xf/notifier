package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin_Success_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	_ = createTestUser(t, deps, "rober0xf", "rober0xf@gmail.com", "password1#!")
	payload := `{"email": "rober0xf@gmail.com", "password": "password1#!"}`

	req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "rober0xf@gmail.com", response["email"])

	assert.NotEmpty(t, response["token"])
	token := response["token"].(string)
	assert.Contains(t, token, ".")
	parts := strings.Split(token, ".")
	assert.Len(t, parts, 3)
}

func TestLogin_MissingFields_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	payload := `{"email": "rober0xf@gmail.com"}`

	req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "email and password are required")
}

func TestLogin_InvalidCredentials_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	_ = createTestUser(t, deps, "rober0xf", "rober0xf@gmail.com", "password1#!")
	payload := `{"email": "rober0xf@gmail.com", "password": "wrongpassword!#"}`

	req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid credentials")
}

func TestLogin_UserNotFound_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	payload := `{"email": "rober0xf@gmail.com", "password": "password1!#"}`

	req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid credentials")
}
