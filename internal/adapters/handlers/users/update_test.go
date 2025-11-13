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
	"golang.org/x/crypto/bcrypt"
)

func TestUpdateUserProfile(t *testing.T) {
	deps := SetupTestDependencies(t)

	username := "rober0xf"
	email := "rober0xf@gmail.com"
	password := "StrongP@ssw0rd!"

	userID := create_test_user(t, deps, username, email, password)
	require.NotZero(t, userID)

	updatePayload := `{
		"username": "updated_rober",
		"email": "rober0xf@gmail.com"
	}`

	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d", userID), strings.NewReader(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "updated_rober", response["username"])
	assert.Equal(t, email, response["email"])

	user, err := deps.userRepo.GetUserByID(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, "updated_rober", user.Username)
}

func TestUpdateUserPassword(t *testing.T) {
	deps := SetupTestDependencies(t)

	username := "rober0xf"
	email := "rober0xf@gmail.com"
	password := "StrongP@ssw0rd!"
	userID := create_test_user(t, deps, username, email, password)
	userIDStr := strconv.Itoa(userID)

	updatePayload := `{
		"password": "StrongP@ssw0rd!2"
	}`

	req := httptest.NewRequest("PUT", "/users/"+userIDStr, strings.NewReader(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	user, err := deps.userRepo.GetUserByID(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, username, user.Username)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("StrongP@ssw0rd!2"))
	assert.NoError(t, err)
}
