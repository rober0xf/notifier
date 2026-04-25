package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Success_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"email": "rober0xf@gmail.com",
		"password": "password1#!"
	}`

	req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "check your email to verify your account", response["message"])
	user := response["user"].(map[string]any)
	assert.NotEmpty(t, user["id"])
	assert.Contains(t, user["email"].(string), "@gmail.com")
	assert.Equal(t, false, user["active"])
}

func TestCreateUser_EmailAlreadyExists_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	originalUser := `{
		"username": "rober0xf",
		"email": "rober@gmail.com",
		"password": "password1#!"
	}`

	req1 := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(originalUser))
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	deps.router.ServeHTTP(w1, req1)
	require.Equal(t, http.StatusCreated, w1.Code)

	duplicateUser := `{
		"username": "rober",
		"email": "rober@gmail.com",
		"password": "password2#!"
	}`

	req2 := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(duplicateUser))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	deps.router.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusConflict, w2.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "email already in use", response.Error)
}

func TestCreateUser_UsernameAlreadyExists_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	originalUser := `{
		"username": "rober0xf",
		"email": "rober@gmail.com",
		"password": "password1#!"
	}`

	req1 := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(originalUser))
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	deps.router.ServeHTTP(w1, req1)
	require.Equal(t, http.StatusCreated, w1.Code)

	duplicateUser := `{
		"username": "rober0xf",
		"email": "rober2@gmail.com",
		"password": "password2#!"
	}`

	req2 := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(duplicateUser))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	deps.router.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusConflict, w2.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "username already in use", response.Error)
}

func TestCreateUser_InvalidEmail_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"email": "invalid@email",
		"password": "password1#!"
	}`

	req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response.Error, "username, email and password are required")
}

func TestCreateUser_WeakPassword_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"email": "valid@gmail.com",
		"password": "aahsdjhsdaaaaaa"
	}`

	req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "password must be stronger", response.Error)
}

func TestCreateUser_MissingFields_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"password": "password1#!"
	}`

	req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "username, email and password are required. password min length 8", response.Error)
}
