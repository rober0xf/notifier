package users

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	deps := SetupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"email": "rober0xf@gmail.com",
		"password": "StrongP@ssw0rd!"
	}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	user, err := deps.userRepo.GetUserByEmail(context.Background(), "rober0xf@gmail.com")
	assert.NoError(t, err)
	assert.Equal(t, "rober0xf", user.Username)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	deps := SetupTestDependencies(t)

	originalUser := `{
		"username": "rober0xf",
		"email": "rober@gmail.com",
		"password": "Strong!P@ssw0rd!"
	}`
	req1 := httptest.NewRequest("POST", "/users", strings.NewReader(originalUser))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	deps.router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	duplicateUser := `{
		"username": "rober",
		"email": "rober@gmail.com",
		"password": "Another$trongP@ss1"
	}`
	req2 := httptest.NewRequest("POST", "/users", strings.NewReader(duplicateUser))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	deps.router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)

	var response map[string]any
	json.Unmarshal(w2.Body.Bytes(), &response)
	assert.Equal(t, "user already exists", response["error"])
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	deps := SetupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"email": "test@disposable-email.com",
		"password": "strongP@ssw0rd!"
	}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "error validating email", response["error"])
}

func TestCreateUser_WeakPassword(t *testing.T) {
	deps := SetupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"email": "valid@gmail.com",
		"password": "aaaaaaa"
	}`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "password must be stronger", response["error"])
}
