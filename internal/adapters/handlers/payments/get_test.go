package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/testutils"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllPayments_Success(t *testing.T) {
	deps := SetupTestDependencies(t)
	userID1, err := testutils.InsertTestUser(context.Background(), deps.db, "rober0xf@gmail.com", "rober0xf")
	userID2, err := testutils.InsertTestUser(context.Background(), deps.db, "rober@gmail.com", "rober0xf")
	assert.NoError(t, err)

	create_test_payment(t, deps, userID1, "Netflix", 25.05, string(domain.Expense), string(domain.Entertainment), "2025-11-22")
	create_test_payment(t, deps, userID2, "HBO", 30, string(domain.Expense), string(domain.Entertainment), "2025-11-22")

	req := httptest.NewRequest("GET", "/payments", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var payments []map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &payments)
	assert.NoError(t, err)

	assert.Len(t, payments, 2)
	assert.Equal(t, "Netflix", payments[0]["name"])
	assert.Equal(t, 25.05, payments[0]["amount"])
	assert.Equal(t, "entertainment", payments[0]["category"])
	assert.Equal(t, "HBO", payments[1]["name"])
	assert.Equal(t, 30.0, payments[1]["amount"])
	assert.Equal(t, "entertainment", payments[1]["category"])
}

func TestGetAllPayments_Empty(t *testing.T) {
	deps := SetupTestDependencies(t)

	req := httptest.NewRequest("GET", "/payments", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var payments []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &payments)
	assert.NoError(t, err)

	assert.Empty(t, payments)
}

func TestGetAllPaymentsFromUser_Success(t *testing.T) {
	deps := SetupTestDependencies(t)

	email := "rober0xf@gmail.com"
	userID, err := testutils.InsertTestUser(context.Background(), deps.db, email, "rober0xf")
	require.NoError(t, err)

	create_test_payment(t, deps, userID, "Netflix", 25.05, string(domain.Expense), string(domain.Entertainment), "2025-11-22")
	create_test_payment(t, deps, userID, "HBO", 30, string(domain.Expense), string(domain.Entertainment), "2025-09-22")
	create_test_payment(t, deps, userID, "Fight Pass", 80, string(domain.Subscription), string(domain.Sports), "2025-10-22")

	req := httptest.NewRequest("GET", "/payments/email?email="+url.QueryEscape(email), nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var payments []map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &payments)
	assert.NoError(t, err)

	assert.Len(t, payments, 3)
	assert.Equal(t, "Netflix", payments[0]["name"])
	assert.Equal(t, 25.05, payments[0]["amount"])
	assert.Equal(t, "entertainment", payments[0]["category"])

	assert.Equal(t, "HBO", payments[1]["name"])
	assert.Equal(t, 30.0, payments[1]["amount"])
	assert.Equal(t, "entertainment", payments[1]["category"])

	assert.Equal(t, "Fight Pass", payments[2]["name"])
	assert.Equal(t, 80.0, payments[2]["amount"])
	assert.Equal(t, "sports", payments[2]["category"])
}

func TestGetAllPaymentsFromUser_Empty(t *testing.T) {
	deps := SetupTestDependencies(t)

	email := "rober0xf@gmail.com"
	_, err := testutils.InsertTestUser(context.Background(), deps.db, email, "rober0xf")
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/payments/email?email="+url.QueryEscape(email), nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var payments []map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &payments)
	assert.NoError(t, err)

	assert.Len(t, payments, 0)
}

func TestGetAllPaymentsFromUser_NotFound(t *testing.T) {
	deps := SetupTestDependencies(t)

	email := "rober0xf@gmail.com"
	userID, err := testutils.InsertTestUser(context.Background(), deps.db, "usernotfound@gmail.com", "rober0xf")
	require.NoError(t, err)
	create_test_payment(t, deps, userID, "Netflix", 25.05, string(domain.Expense), string(domain.Entertainment), "2025-11-22")

	req := httptest.NewRequest("GET", "/payments/email?email="+url.QueryEscape(email), nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "user not found", errorResponse["error"])
}

func TestGetPaymentByID_Sucess(t *testing.T) {
	deps := SetupTestDependencies(t)

	userID, err := testutils.InsertTestUser(context.Background(), deps.db, "rober0xf@gmail.com", "rober0xf")
	require.NoError(t, err)
	paymentID := create_test_payment(t, deps, userID, "Netflix", 25.05, string(domain.Expense), string(domain.Entertainment), "2025-11-22")

	req := httptest.NewRequest("GET", fmt.Sprintf("/payments/%d", paymentID), nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPaymentByID_NotFound(t *testing.T) {
	deps := SetupTestDependencies(t)

	req := httptest.NewRequest("GET", "/payments/2", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
