package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMyPayments_Success_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)
	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)

	_, err = insertTestPayment(context.Background(), paymentDeps.db, userID, TestPayment{
		Name:      "claude",
		Amount:    15.99,
		Type:      "subscription",
		Category:  "work",
		Date:      "2026-01-03",
		Paid:      false,
		Recurrent: false,
	})
	require.NoError(t, err)
	_, err = insertTestPayment(context.Background(), paymentDeps.db, userID, TestPayment{
		Name:      "mac",
		Amount:    1700,
		Type:      "expense",
		Category:  "electronics",
		Date:      "2026-01-21",
		Paid:      true,
		Recurrent: false,
	})
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/v1/auth/payments/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response []dto.PaymentResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)

	assert.Equal(t, "claude", response[0].Name)
	assert.Equal(t, "mac", response[1].Name)
	assert.False(t, response[0].Paid)
	assert.True(t, response[1].Paid)
}

func TestGetMyPayments_ReturnsOnlyOwnPayments_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)
	token2 := getAuthToken(t, userDeps)
	userID1, err := getUserIDFromToken(token)
	require.NoError(t, err)
	userID2, err := getUserIDFromToken(token2)
	require.NoError(t, err)

	_, err = insertTestPayment(context.Background(), paymentDeps.db, userID1, TestPayment{
		Name:      "claude",
		Amount:    15.99,
		Type:      "subscription",
		Category:  "work",
		Date:      "2026-01-03",
		Paid:      false,
		Recurrent: false,
	})
	require.NoError(t, err)
	_, err = insertTestPayment(context.Background(), paymentDeps.db, userID1, TestPayment{
		Name:      "mac",
		Amount:    1700,
		Type:      "expense",
		Category:  "work",
		Date:      "2026-01-03",
		Paid:      true,
		Recurrent: false,
	})
	require.NoError(t, err)
	_, err = insertTestPayment(context.Background(), paymentDeps.db, userID2, TestPayment{
		Name:      "tv",
		Amount:    900,
		Type:      "expense",
		Category:  "electronics",
		Date:      "2026-05-15",
		Paid:      true,
		Recurrent: false,
	})
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/v1/auth/payments/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response []dto.PaymentResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response, 2)
	assert.Equal(t, "claude", response[0].Name)
	assert.Equal(t, "mac", response[1].Name)
}

func TestGetMyPayments_EmptyList_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/auth/payments/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response []dto.PaymentResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 0)
}

func TestGetMyPayments_Unauthorized_Integration(t *testing.T) {
	_, paymentDeps := setupTestDependencies(t)

	req := httptest.NewRequest("GET", "/v1/auth/payments/me", nil)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "no token provided", response.Error)
}
