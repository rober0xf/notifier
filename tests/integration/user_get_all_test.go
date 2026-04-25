package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllUsers_Success_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAdminToken(t, userDeps)

	getAuthToken(t, userDeps)
	getAuthToken(t, userDeps)
	getAuthToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response []dto.UserPayload
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(response), 3)
}

func TestGetAllUsers_Empty_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAdminToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response []dto.UserPayload
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(response), 1)
}

func TestGetAllUsers_Unauthorized_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)

	req := httptest.NewRequest("GET", "/v1/admin/users", nil)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "no token provided", response.Error)
}

func TestGetAllUsers_Forbidden_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "forbidden", response.Error)
}
