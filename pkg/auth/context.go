package auth

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(c *gin.Context) (int, error) {
	userIDcontext, exists := c.Get("user_id")
	if !exists {
		slog.InfoContext(c.Request.Context(), "failed to get user_id from context")
		return 0, ErrUserIDNotFound
	}

	userID, ok := userIDcontext.(int)
	if !ok {
		slog.InfoContext(c.Request.Context(), "failed to cast user_id from context")
		return 0, ErrInvalidUserID
	}

	return userID, nil
}
