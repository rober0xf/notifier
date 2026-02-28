package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Success_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"email": "rober0xf@gmail.com",
		"password": "password1#!"
	}`

	req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	require.Equal(t, "rober0xf@gmail.com", response["email"])
	assert.NotEmpty(t, response["token"])

	user, err := deps.userRepo.GetUserByEmail(context.Background(), "rober0xf@gmail.com")
	require.NoError(t, err)
	require.Equal(t, "rober0xf", user.Username)
	assert.False(t, user.Active)
}

func TestCreateUser_AlreadyExists_Integration(t *testing.T) {
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

	var response map[string]any
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user already exists", response["error"])
}

func TestCreateUser_InvalidEmail_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"email": "invalid-email.com",
		"password": "password1#!"
	}`

	req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_DisposableEmail_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"email": "ignoremail.com",
		"password": "password1#!"
	}`

	req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_WeakPassword_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"email": "valid@gmail.com",
		"password": "aaaaaaa"
	}`

	req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_MissingFields_Integration(t *testing.T) {
	deps, _ := setupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"password": "password1#!"
	}`

	req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
