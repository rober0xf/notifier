package integration

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/stretchr/testify/require"
)

func getAuthToken(t *testing.T, deps *TestUserDependencies) string {
	t.Helper()

	ctx := context.Background()
	id := uuid.New().String()

	email := fmt.Sprintf("test-%s@example.com", id)
	username := fmt.Sprintf("test-%s", id)
	password := fmt.Sprintf("test-%s", id)

	userID, err := insertTestUser(ctx, deps.db, email, username, password)
	require.NoError(t, err)

	token, err := deps.jwtGen.Generate(userID, email, entity.RoleUser)
	require.NoError(t, err)

	return token
}

func getAdminToken(t *testing.T, deps *TestUserDependencies) string {
	t.Helper()

	ctx := context.Background()
	id := uuid.New().String()

	email := fmt.Sprintf("test-%s@example.com", id)
	username := fmt.Sprintf("test-%s", id)
	password := fmt.Sprintf("test-%s", id)

	userID, err := insertTestUser(ctx, deps.db, email, username, password)
	require.NoError(t, err)

	token, err := deps.jwtGen.Generate(userID, email, entity.RoleAdmin)
	require.NoError(t, err)

	return token
}

func getUserIDFromToken(token string) (int, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, err
	}

	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return 0, err
	}

	userID := int(claims["user_id"].(float64))
	return userID, nil
}

func extractEmailFromToken(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", err
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("email claim not found")
	}

	return email, nil
}
