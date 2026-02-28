package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeletePayment_Success_Integration(t *testing.T) {
	_, deps := setupTestDependencies(t)
	token := getAuthToken(t, deps.router)

	paymentID := createTestPayment(t, deps, token, "claude", 50.99, "subscription", "work", "2026-01-03", false, false)
	paymentIDStr := strconv.Itoa(paymentID)

	req := httptest.NewRequest("DELETE", "/v1/auth/payments/"+paymentIDStr, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)

	_, err := deps.paymentRepo.GetPaymentByID(context.Background(), paymentID)
	assert.Equal(t, "resource not found", err.Error())
}

func TestDeletePayment_Unauthorized_Integration(t *testing.T) {
	_, deps := setupTestDependencies(t)

	req := httptest.NewRequest("DELETE", "/v1/auth/payments/1", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeletePayment_InvalidID_Integration(t *testing.T) {
	_, deps := setupTestDependencies(t)
	token := getAuthToken(t, deps.router)

	req := httptest.NewRequest("DELETE", "/v1/auth/payments/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid payment id", response["error"])
}

func TestDeletePayment_NegativeID_Integration(t *testing.T) {
	_, deps := setupTestDependencies(t)
	token := getAuthToken(t, deps.router)

	req := httptest.NewRequest("DELETE", "/v1/auth/payments/-1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "id must be positive", response["error"])
}

func TestDeletePayment_NotFound_Integration(t *testing.T) {
	_, deps := setupTestDependencies(t)
	token := getAuthToken(t, deps.router)

	req := httptest.NewRequest("DELETE", "/v1/auth/payments/9999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "payment not found", response["error"])
}

func TestDeletePayment_Forbidden_Integration(t *testing.T) {
	_, deps := setupTestDependencies(t)
	token := getAuthToken(t, deps.router)

	paymentID := createTestPayment(t, deps, token, "claude", 50.99, "subscription", "work", "2026-01-03", false, false)
	paymentIDStr := strconv.Itoa(paymentID)

	otherToken := getAuthTokenWithCredentials(t, deps.router, "other", "other@gmail.com", "password1#!")
	req := httptest.NewRequest("DELETE", "/v1/auth/payments/"+paymentIDStr, nil)
	req.Header.Set("Authorization", "Bearer "+otherToken)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "cannot delete payments from other user", response["error"])
}
