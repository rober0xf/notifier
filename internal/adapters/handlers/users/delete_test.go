package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser_Success(t *testing.T) {
	deps := SetupTestDependencies(t)

	userID := create_test_user(t, deps, "rober", "rober0xf@gmail.com", "securePassword!")
	userIDStr := strconv.Itoa(userID)

	req := httptest.NewRequest("DELETE", "/users/"+userIDStr, nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	_, err := deps.userRepo.GetUserByEmail(context.Background(), "rober0xf@gmail.com")
	assert.ErrorIs(t, err, dto.ErrNotFound)
}

func TestDeleteUser_NotFound(t *testing.T) {
	deps := SetupTestDependencies(t)

	nonExistingID := strconv.Itoa(999)
	req := httptest.NewRequest("DELETE", "/users/"+nonExistingID, nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

}
