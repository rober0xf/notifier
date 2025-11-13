package payments

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/adapters/testutils"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestDeletePayment_Success(t *testing.T) {
	deps := SetupTestDependencies(t)
	userID, err := testutils.InsertTestUser(context.Background(), deps.db, "rober0xf@gmail.com", "rober0xf")
	assert.NoError(t, err)

	paymentID := create_test_payment(t, deps, userID, "Work", 1200, string(domain.Income), string(domain.Work), "2025-12-02")

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/payments/%d", paymentID), nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	_, err = deps.paymentRepo.GetPaymentByID(context.Background(), paymentID)
	assert.ErrorIs(t, err, dto.ErrNotFound)
}

func TestDeletePayment_NotFound(t *testing.T) {
	deps := SetupTestDependencies(t)

	req := httptest.NewRequest("DELETE", "/payments/999", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
