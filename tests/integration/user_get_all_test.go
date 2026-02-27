package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllUsers_Success_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	_ = createTestUser(t, deps, "user1", "usermail1@gmail.com", "password1!#")
	_ = createTestUser(t, deps, "user2", "usermail2@gmail.com", "password2!#")

	req := httptest.NewRequest("GET", "/v1/users", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var users []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &users)
	require.NoError(t, err)

	require.Equal(t, 2, len(users))
	require.Equal(t, "user1", users[0]["username"])
	require.Equal(t, "user2", users[1]["username"])
}

func TestGetAllUsers_Empty_Integration(t *testing.T) {
	deps := setupTestUserDependencies(t)

	req := httptest.NewRequest("GET", "/v1/users", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var users []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &users)
	require.NoError(t, err)
	require.Empty(t, users)
}
