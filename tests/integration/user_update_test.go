package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUpdateUser_Success_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	userID := createTestUser(t, deps, "rober0xf", "rober0xf@gmail.com", "password1#!")
	userIDStr := strconv.Itoa(userID)

	input := `{"username": "rober"}`

	req := httptest.NewRequest("PATCH", "/v1/users/"+userIDStr, strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "rober", response["username"])
	assert.Equal(t, "rober0xf@gmail.com", response["email"])
}

func TestUpdateUser_InvalidID_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	req := httptest.NewRequest("PATCH", "/v1/users/aa", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid id")
}

func TestUpdateUser_MalformedJSON_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	userID := createTestUser(t, deps, "rober0xf", "rober0xf@gmail.com", "password1#!")
	userIDStr := strconv.Itoa(userID)

	input := `{invalid}`

	req := httptest.NewRequest("PATCH", "/v1/users/"+userIDStr, strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid request")
}

func TestUpdateUser_EmptyRequest_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	userID := createTestUser(t, deps, "rober0xf", "rober0xf@gmail.com", "password1#!")
	userIDStr := strconv.Itoa(userID)

	input := `{}`

	req := httptest.NewRequest("PATCH", "/v1/users/"+userIDStr, strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "at least one field is required")
}

func TestUpdateUser_NotFound_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	input := `{"username": "rober"}`

	req := httptest.NewRequest("PATCH", "/v1/users/9999", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "user not found")
}

func TestUpdateUser_NegativeID_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	input := `{"username": "rober"}`

	req := httptest.NewRequest("PATCH", "/v1/users/-1", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "id must be positive")
}

func TestUpdateUser_InvalidEmail_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	userID := createTestUser(t, deps, "rober0xf", "rober0xf@gmail.com", "password1#!")
	userIDStr := strconv.Itoa(userID)

	input := `{"email": "rober.com"}`

	req := httptest.NewRequest("PATCH", "/v1/users/"+userIDStr, strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid email format")
}

func TestUpdateUser_Password_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	userID := createTestUser(t, deps, "rober0xf", "rober0xf@gmail.com", "password1!#")
	userIDStr := strconv.Itoa(userID)

	input := `{"password": "password2!#"}`

	req := httptest.NewRequest("PATCH", "/v1/users/"+userIDStr, strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	user, err := deps.userRepo.GetUserByID(context.Background(), userID)
	require.NoError(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password2!#"))
	assert.NoError(t, err)
}
