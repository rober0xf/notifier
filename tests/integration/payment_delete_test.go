package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeletePayment_Success_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)
	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)

	paymentID, err := insertTestPayment(context.Background(), paymentDeps.db, userID, TestPayment{
		Name:      "claude",
		Amount:    15.99,
		Type:      "subscription",
		Category:  "work",
		Date:      "2026-01-03",
		Paid:      false,
		Recurrent: false,
	})
	require.NoError(t, err)

	req := httptest.NewRequest("DELETE", "/v1/auth/payments/"+strconv.Itoa(paymentID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)

	_, err = paymentDeps.paymentRepo.GetPaymentByID(context.Background(), paymentID)
	assert.ErrorIs(t, err, errors.ErrNotFound)
}

func TestDeletePayment_Unauthorized_Integration(t *testing.T) {
	_, paymentDeps := setupTestDependencies(t)

	req := httptest.NewRequest("DELETE", "/v1/auth/payments/1", nil)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "no token provided", response.Error)
}

func TestDeletePayment_InvalidID_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	req := httptest.NewRequest("DELETE", "/v1/auth/payments/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid payment id", response.Error)
}

func TestDeletePayment_NegativeID_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	req := httptest.NewRequest("DELETE", "/v1/auth/payments/-1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid payment id", response.Error)
}

func TestDeletePayment_NotFound_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	req := httptest.NewRequest("DELETE", "/v1/auth/payments/9999", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "payment not found", response.Error)
}

func TestDeletePayment_Forbidden_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)
	token2 := getAuthToken(t, userDeps)
	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)

	paymentID, err := insertTestPayment(context.Background(), paymentDeps.db, userID, TestPayment{
		Name:      "claude",
		Amount:    15.99,
		Type:      "subscription",
		Category:  "work",
		Date:      "2026-01-03",
		Paid:      false,
		Recurrent: false,
	})
	require.NoError(t, err)

	req := httptest.NewRequest("DELETE", "/v1/auth/payments/"+strconv.Itoa(paymentID), nil)
	req.Header.Set("Authorization", "Bearer "+token2)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)

	var response dto.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "cannot delete payment from other user", response.Error)
}
