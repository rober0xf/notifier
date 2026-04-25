package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePayment_Success(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	payload := `{
			"name": "claude",
			"amount": 50.99,
			"type": "subscription",
			"category": "work",
			"date": "2026-01-03",
			"paid": false
		}`

	req := httptest.NewRequest("POST", "/v1/auth/payments", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.NotEmpty(t, w.Body.Bytes())

	var resp dto.PaymentResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	require.Equal(t, "claude", resp.Name)
	require.Equal(t, 50.99, resp.Amount)
	require.Equal(t, entity.TransactionTypeSubscription, resp.Type)
	require.Equal(t, entity.CategoryTypeWork, resp.Category)
	require.False(t, resp.Paid)

	require.Nil(t, resp.DueDate)
	require.Nil(t, resp.PaidAt)
	require.Nil(t, resp.Frequency)
	require.Nil(t, resp.ReceiptURL)

	payment, err := paymentDeps.paymentRepo.GetPaymentByID(context.Background(), int(resp.ID))
	require.NoError(t, err)

	require.Equal(t, "claude", payment.Name)
	require.Equal(t, 50.99, payment.Amount)
	require.Equal(t, entity.CategoryTypeWork, payment.Category)
	require.Equal(t, entity.TransactionTypeSubscription, payment.Type)
	require.Nil(t, payment.Frequency)
	assert.False(t, payment.Paid)
}

func TestCreatePayment_MissingName_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	payload := `{
        "amount": 50.99,
        "type": "subscription",
        "category": "work",
        "date": "2026-01-03",
        "paid": false,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/auth/payments", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "name is required", response.Error)
}

func TestCreatePayment_InvalidType_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "invalid",
        "category": "work",
        "date": "2026-01-03",
        "paid": false,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/auth/payments", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response.Error, "invalid transaction type")
}

func TestCreatePayment_InvalidCategory_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "subscription",
        "category": "invalid",
        "date": "2026-01-03",
        "paid": false,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/auth/payments", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response.Error, "invalid category type")
}

func TestCreatePayment_InvalidDate_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "subscription",
        "category": "work",
        "date": "03-01-2026",
        "paid": false,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/auth/payments", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "date must be in YYYY-MM-DD format", response.Error)
}

func TestCreatePayment_Unauthorized_Integration(t *testing.T) {
	_, paymentDeps := setupTestDependencies(t)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "subscription",
        "category": "work",
        "date": "2026-01-03",
        "paid": false,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/auth/payments", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "no token provided", response.Error)
}

func TestCreatePayment_RecurrentWithoutFrequency_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "subscription",
        "category": "work",
        "date": "2026-01-03",
        "paid": false,
        "recurrent": true
    }`

	req := httptest.NewRequest("POST", "/v1/auth/payments", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response.Error, "invalid frequency type")
}

func TestCreatePayment_PaidWithoutPaidAt_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "subscription",
        "category": "work",
        "date": "2026-01-03",
        "paid": true,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/auth/payments", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var response dto.PaymentResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotNil(t, response.PaidAt)
}
