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

// change when route is auth
// func TestDeleteUser_Success_OwnAccount_Integration(t *testing.T) {
// 	userDeps, _ := setupTestDependencies(t)
// 	token := getAuthToken(t, userDeps)
//
// 	userID, _ := extractUserIDFromToken(token)
// 	userIDStr := strconv.Itoa(userID)
//
// 	req := httptest.NewRequest("DELETE", "/v1/admin/users/"+userIDStr, nil)
// 	req.Header.Set("Authorization", "Bearer "+token)
//
// 	w := httptest.NewRecorder()
// 	userDeps.router.ServeHTTP(w, req)
// 	require.Equal(t, http.StatusNoContent, w.Code)
// }

func TestDeleteUser_Success_Admin_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAdminToken(t, userDeps)

	userID, _ := getUserIDFromToken(token)
	userIDStr := strconv.Itoa(userID)

	req := httptest.NewRequest("DELETE", "/v1/admin/users/"+userIDStr, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteUser_InvalidID_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAdminToken(t, userDeps)

	req := httptest.NewRequest("DELETE", "/v1/admin/users/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid id", response.Error)
}

func TestDeleteUser_Unauthorized_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)

	req := httptest.NewRequest("DELETE", "/v1/admin/users/1", nil)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "no token provided", response.Error)
}

func TestDeleteUser_Forbidden_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)
	token2 := getAuthToken(t, userDeps)
	userID, _ := getUserIDFromToken(token)
	userIDStr := strconv.Itoa(userID)

	req := httptest.NewRequest("DELETE", "/v1/admin/users/"+userIDStr, nil)
	req.Header.Set("Authorization", "Bearer "+token2)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "forbidden", response.Error)
}

func TestDeleteUser_NotFound_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	adminToken := getAdminToken(t, userDeps)

	req := httptest.NewRequest("DELETE", "/v1/admin/users/999999", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user not found", response.Error)
}
