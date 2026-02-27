package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePayment_Success(t *testing.T) {
	deps := setupTestPaymentDependencies(t)
	token := getAuthToken(t, deps)

	payload := `{
			"name": "claude",
			"amount": 50.99,
			"type": "subscription",
			"category": "work",
			"date": "2026-01-03",
			"paid": false,
			"recurrent": false
		}`

	req := httptest.NewRequest("POST", "/v1/payments/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer"+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)

	require.NoError(t, err)
	require.Equal(t, "claude", response["name"])

	id := int(response["id"].(float64))
	payment, err := deps.paymentRepo.GetPaymentByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, "work", payment.Category)
	assert.False(t, payment.Paid)
}

func TestCreatePayment_MissingName_Integration(t *testing.T) {
	deps := setupTestPaymentDependencies(t)
	token := getAuthToken(t, deps)

	payload := `{
        "amount": 50.99,
        "type": "subscription",
        "category": "work",
        "date": "2026-01-03",
        "paid": false,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/payments/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "name is required", response["error"])
}

func TestCreatePayment_InvalidType_Integration(t *testing.T) {
	deps := setupTestPaymentDependencies(t)
	token := getAuthToken(t, deps)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "invalid",
        "category": "work",
        "date": "2026-01-03",
        "paid": false,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/payments/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "type must be: expense, income, or subscription", response["error"])
}

func TestCreatePayment_InvalidCategory_Integration(t *testing.T) {
	deps := setupTestPaymentDependencies(t)
	token := getAuthToken(t, deps)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "subscription",
        "category": "invalid",
        "date": "2026-01-03",
        "paid": false,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/payments/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "category must be: electronics, entertainment, education, clothing, work, or sports", response["error"])
}

func TestCreatePayment_InvalidDate_Integration(t *testing.T) {
	deps := setupTestPaymentDependencies(t)
	token := getAuthToken(t, deps)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "subscription",
        "category": "work",
        "date": "03-01-2026",
        "paid": false,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/payments/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "date must be in YYYY-MM-DD format", response["error"])
}

func TestCreatePayment_Unauthorized_Integration(t *testing.T) {
	deps := setupTestPaymentDependencies(t)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "subscription",
        "category": "work",
        "date": "2026-01-03",
        "paid": false,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/payments/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "unauthorized", response["error"])
}

func TestCreatePayment_AlreadyExists_Integration(t *testing.T) {
	deps := setupTestPaymentDependencies(t)
	token := getAuthToken(t, deps)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "subscription",
        "category": "work",
        "date": "2026-01-03",
        "paid": false,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/payments/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	req = httptest.NewRequest("POST", "/v1/payments/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusConflict, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "payment already exists", response["error"])
}

func TestCreatePayment_RecurrentWithoutFrequency_Integration(t *testing.T) {
	deps := setupTestPaymentDependencies(t)
	token := getAuthToken(t, deps)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "subscription",
        "category": "work",
        "date": "2026-01-03",
        "paid": false,
        "recurrent": true
    }`

	req := httptest.NewRequest("POST", "/v1/payments/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "frequency required for recurrent payments", response["error"])
}

func TestCreatePayment_PaidWithoutPaidAt_Integration(t *testing.T) {
	deps := setupTestPaymentDependencies(t)
	token := getAuthToken(t, deps)

	payload := `{
        "name": "claude",
        "amount": 50.99,
        "type": "subscription",
        "category": "work",
        "date": "2026-01-03",
        "paid": true,
        "recurrent": false
    }`

	req := httptest.NewRequest("POST", "/v1/payments/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "paid_at is required when paid is true", response["error"])
}
