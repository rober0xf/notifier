package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserByID_Success_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAdminToken(t, userDeps)
	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/v1/admin/users/"+strconv.Itoa(userID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response dto.UserPayload
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, userID, response.ID)
}

func TestGetUserByID_InvalidID_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	adminToken := getAdminToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/users/abb", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid id", response.Error)
}

func TestGetUserByID_NotFound_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	adminToken := getAdminToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/users/999999", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user not found", response.Error)
}

func TestGetUserByID_Unauthorized_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)
	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/v1/admin/users/"+strconv.Itoa(userID), nil)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "no token provided", response.Error)
}

func TestGetUserByID_Forbidden_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)
	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/v1/admin/users/"+strconv.Itoa(userID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)

	var response dto.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "forbidden", response.Error)
}
