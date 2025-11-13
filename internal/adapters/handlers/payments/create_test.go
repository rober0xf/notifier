package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/testutils"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper
func create_test_payment(t *testing.T, deps *TestDependencies, userID int, name string, amount float64, txType, category, date string) int {
	payload := fmt.Sprintf(`{
		"name": "%s",
		"amount": %f,
		"type": "%s",
		"category": "%s",
		"date": "%s",
		"paid": false,
		"recurrent": false
	}`, name, amount, txType, category, date)

	req := httptest.NewRequest("POST", "/payments", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer testtoken")

	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Logf("unexpected status: %d, body: %s", w.Code, w.Body.String())
	}
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return int(response["id"].(float64))
}

func TestCreatePayment_Success(t *testing.T) {
	deps := SetupTestDependencies(t)
	userID, err := testutils.InsertTestUser(context.Background(), deps.db, "rober0xf@gmail.com", "rober0xf")
	assert.NoError(t, err)

	paymentID := create_test_payment(t, deps, userID, "Groceries", 50.0, string(domain.Expense), string(domain.Clothing), "2025-11-11")
	assert.NotZero(t, paymentID)
}

func TestCreatePayment_MissingFields(t *testing.T) {
	deps := SetupTestDependencies(t)

	payload := `{"name": "Test"}`
	req := httptest.NewRequest("POST", "/payments", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePayment_InvalidJSON(t *testing.T) {
	deps := SetupTestDependencies(t)

	req := httptest.NewRequest("POST", "/payments", strings.NewReader(`{"name":123}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
