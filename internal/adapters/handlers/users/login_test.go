package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin_Success(t *testing.T) {
	deps := SetupTestDependencies(t)

	email := "rober0xf@gmail.com"
	password := "StrongP@ssw0rd!"
	create_test_user(t, deps, "rober0xf", email, password)
	loginPayload := fmt.Sprintf(`{
		"email": "%s",
		"password": "%s"
	}`, email, password)

	req := httptest.NewRequest("POST", "/users/login", strings.NewReader(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	user := response["user"].(map[string]any)
	assert.NotEmpty(t, user["token"])
	assert.Equal(t, email, user["email"])
	assert.Equal(t, float64(1), user["id"])

	token := user["token"].(string)
	assert.Contains(t, token, ".")
	parts := strings.Split(token, ".")
	assert.Len(t, parts, 3)
}

func TestLogin_WrongPassword(t *testing.T) {
	deps := SetupTestDependencies(t)

	email := "rober0xf@gmail.com"
	password := "StrongP@ssw0rd!"
	create_test_user(t, deps, "rober0xf", email, password)
	loginPayload := fmt.Sprintf(`{
		"email": "%s",
		"password": "WrongPassword123!"
	}`, email)

	req := httptest.NewRequest("POST", "/users/login", strings.NewReader(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "incorrect password", response["error"])
}

func TestLogin_UserNotFound(t *testing.T) {
	deps := SetupTestDependencies(t)

	loginPayload := `{
		"email": "doesnotexists@gmail.com",
		"password": "SecurePassword123!"
	}`
	req := httptest.NewRequest("POST", "/users/login", strings.NewReader(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user not found", response["error"])
}

func TestLogin_MissingCredentials(t *testing.T) {
	deps := SetupTestDependencies(t)

	tests := []struct {
		name    string
		payload string
	}{
		{
			name:    "missing email",
			payload: `{"password": "SomePassword123!"}`,
		},
		{
			name:    "missing password",
			payload: `{"email": "test@gmail.com"}`,
		},
		{
			name:    "empty payload",
			payload: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/users/login", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			deps.router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]any
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, "email and password are required", response["error"])
		})
	}
}
