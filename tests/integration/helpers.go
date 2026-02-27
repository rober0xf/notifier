package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T, deps *TestUserDependencies, username, email, password string) int {
	t.Helper()

	payload := fmt.Sprintf(`{
        "username": "%s",
        "email": "%s",
        "password": "%s"
    }`, username, email, password)

	req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	user, err := deps.userRepo.GetUserByEmail(context.Background(), email)
	require.NoError(t, err)

	return user.ID
}

func createTestPayment(
	t *testing.T,
	deps *TestPaymentDependencies,
	token string,
	name string,
	amount float64,
	payment_type string,
	category string,
	date string,
	paid bool, recurrent bool) int {
	t.Helper()

	payload := fmt.Sprintf(`{
        "user_id": %d,
        "amount": %f,
        "payment_type": "%s",
        "category": "%s",
        "date": "%s",
        "paid": %t,
        "recurrent": %t
    }`, amount, payment_type, category, date, paid, recurrent)

	req := httptest.NewRequest("POST", "/v1/payments/create", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var response map[string]any
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	return int(response["id"].(float64))
}

func getAuthToken(t *testing.T, deps *TestPaymentDependencies) string {
	return getAuthTokenWithCredentials(t, deps, "rober0xf", "rober0xf@gmail.com", "password1!#")
}

func getAuthTokenWithCredentials(t *testing.T, deps *TestPaymentDependencies, username, email, password string) string {
	t.Helper()

	payload := fmt.Sprintf(`{"username": "%s", "email": "%s", "password": "%s"}`, username, email, password)
	req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response["token"].(string)
}
