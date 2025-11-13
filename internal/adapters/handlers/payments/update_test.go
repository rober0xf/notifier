package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/testutils"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdatePayment_Success(t *testing.T) {
	deps := SetupTestDependencies(t)

	userID, err := testutils.InsertTestUser(context.Background(), deps.db, "rober0xf@gmail.com", "rober0xf")
	require.NoError(t, err)

	paymentID := create_test_payment(t, deps, userID, "Groceries", 50.0, string(domain.Income), string(domain.Entertainment), "2025-11-11")
	require.NotZero(t, paymentID)

	_name := "spotify"
	_type := "subscription"
	updatePayment := domain.UpdatePayment{
		Name: &_name,
		Type: (*domain.TransactionType)(&_type),
	}
	body, err := json.Marshal(updatePayment)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/payments/%d", paymentID), bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer testtoken")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	payment, err := deps.paymentService.Get(context.Background(), paymentID)
	require.NoError(t, err)
	assert.Equal(t, "spotify", response["name"])
	assert.Equal(t, "subscription", response["type"])
	assert.Equal(t, 50.0, payment.Amount)
}

func TestUpdatePayment_NotFound(t *testing.T) {
	deps := SetupTestDependencies(t)

	userID, err := testutils.InsertTestUser(context.Background(), deps.db, "rober0xf@gmail.com", "rober0xf")
	require.NoError(t, err)

	paymentID := create_test_payment(t, deps, userID, "Groceries", 50.0, string(domain.Income), string(domain.Entertainment), "2025-11-11")
	require.NotZero(t, paymentID)

	_name := "spotify"
	_type := "subscription"
	updatePayment := domain.UpdatePayment{
		Name: &_name,
		Type: (*domain.TransactionType)(&_type),
	}
	body, err := json.Marshal(updatePayment)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "/payments/99", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer testtoken")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "payment not found", response["error"])

	existing, err := deps.paymentService.Get(context.Background(), paymentID)
	require.NoError(t, err)
	assert.Equal(t, "Groceries", existing.Name)
}

func TestUpdatePayment_Unauthorized(t *testing.T) {
	deps := SetupTestDependencies(t)

	userID, err := testutils.InsertTestUser(context.Background(), deps.db, "rober0xf@gmail.com", "rober0xf")
	require.NoError(t, err)

	paymentID := create_test_payment(t, deps, userID, "Groceries", 50.0, string(domain.Income), string(domain.Entertainment), "2025-11-11")
	require.NotZero(t, paymentID)

	_name := "spotify"
	_type := "subscription"
	updatePayment := domain.UpdatePayment{
		Name: &_name,
		Type: (*domain.TransactionType)(&_type),
	}
	body, err := json.Marshal(updatePayment)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "/payments/99", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "unauthorized", resp["error"])
}
