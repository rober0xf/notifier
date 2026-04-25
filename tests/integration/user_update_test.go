package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateUser_Success_OwnAccount_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)
	payload := `{
		"username": "changed"
	}`

	req := httptest.NewRequest("PUT", "/v1/auth/users/"+strconv.Itoa(userID), strings.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response dto.UserPayload
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, userID, response.ID)
	assert.Equal(t, "changed", response.Username)
}

func TestUpdateUser_Success_Admin_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)
	adminToken := getAdminToken(t, userDeps)

	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)
	payload := `{
		"username": "changed"
	}`

	req := httptest.NewRequest("PUT", "/v1/auth/users/"+strconv.Itoa(userID), strings.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response dto.UserPayload
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "changed", response.Username)
}

func TestUpdateUser_InvalidID_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)
	payload := `{
		"email": "changed@gmail.com"
	}`

	req := httptest.NewRequest("PUT", "/v1/auth/users/abb", strings.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid id", response.Error)
}

func TestUpdateUser_InvalidBody_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)
	payload := `{
		"position": "swe"
	}`

	req := httptest.NewRequest("PUT", "/v1/auth/users/"+strconv.Itoa(userID), strings.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "at least one field is required", response.Error)
}

func TestUpdateUser_Unauthorized_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)

	payload := `{
		"username": "swe"
	}`

	req := httptest.NewRequest("PUT", "/v1/auth/users/1", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "no token provided", response.Error)
}

func TestUpdateUser_Forbidden_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)
	token2 := getAuthToken(t, userDeps)

	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)
	userIDStr := strconv.Itoa(userID)
	payload := `{
		"username": "swe"
	}`

	req := httptest.NewRequest("PUT", "/v1/auth/users/"+userIDStr, strings.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+token2)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)

	var response dto.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "cannot change other user data", response.Error)
}

func TestUpdateUser_NotFound_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	adminToken := getAdminToken(t, userDeps)

	payload := `{
		"username": "swe"
	}`

	req := httptest.NewRequest("PUT", "/v1/auth/users/999999", strings.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user not found", response.Error)
}

func TestUpdateUser_EmailAlreadyExists_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)
	token2 := getAuthToken(t, userDeps)

	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)
	email2, err := extractEmailFromToken(token2)
	require.NoError(t, err)

	body := fmt.Sprintf(`{"email":"%s"}`, email2)
	req := httptest.NewRequest("PUT", "/v1/auth/users/"+strconv.Itoa(userID), strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusConflict, w.Code)

	var response dto.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "email already in use", response.Error)
}
