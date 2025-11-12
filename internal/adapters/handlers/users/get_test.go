package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper
func create_test_user(t *testing.T, deps *TestDependencies, username, email, password string) int {
	payload := fmt.Sprintf(`{"username":"%s","email":"%s","password":"%s"}`, username, email, password)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	// debug
	if w.Code != http.StatusCreated {
		t.Logf("error here, status: %d, body: %s", w.Code, w.Body.String())
	}

	require.Equal(t, http.StatusCreated, w.Code)

	user, err := deps.userRepo.GetUserByEmail(context.Background(), email)
	require.NoError(t, err)

	return user.ID
}

func TestGetUserByID_Sucess(t *testing.T) {
	deps := SetupTestDependencies(t)
	userID := create_test_user(t, deps, "rober0xf", "rober0xf@gmail.com", "securePassword!")

	userIDStr := strconv.Itoa(userID)
	req := httptest.NewRequest("GET", "/users/"+userIDStr, nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserByID_NotFound(t *testing.T) {
	deps := SetupTestDependencies(t)
	nonExistentUserID := 9999

	userIDStr := strconv.Itoa(nonExistentUserID)
	req := httptest.NewRequest("GET", "/users/"+userIDStr, nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetByEmail_Empty(t *testing.T) {
	deps := SetupTestDependencies(t)

	req := httptest.NewRequest("GET", "/users/email", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetByEmail_Success(t *testing.T) {
	deps := SetupTestDependencies(t)
	email := "rober0xf@gmail.com"

	create_test_user(t, deps, "rober0xf", email, "securePassword!")
	req := httptest.NewRequest("GET", "/users/email/"+email, nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "rober0xf", response["username"])
	assert.Equal(t, email, response["email"])
	assert.NotContains(t, response, "password")
}

func TestGetByEmail_NotFound(t *testing.T) {
	deps := SetupTestDependencies(t)

	nonExistentEmail := "doesnotexists@example.com"
	req := httptest.NewRequest("GET", "/users/email/"+nonExistentEmail, nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAllUsers_Success(t *testing.T) {
	deps := SetupTestDependencies(t)
	create_test_user(t, deps, "user1", "usermail1@gmail.com", "securePassword1#")
	create_test_user(t, deps, "user2", "usermail2@gmail.com", "securePassword2#")

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var users []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &users)
	assert.NoError(t, err)

	assert.Equal(t, len(users), 2)
	assert.Equal(t, "user1", users[0]["username"])
	assert.Equal(t, "user2", users[1]["username"])
}

func TestGetAllUsers_Empty(t *testing.T) {
	deps := SetupTestDependencies(t)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var users []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &users)
	assert.NoError(t, err)

	assert.Empty(t, users)
}
