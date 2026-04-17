package auth

import (
	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(c *gin.Context) (int, error) {
	userIDcontext, exists := c.Get("user_id")
	if !exists {
		return 0, ErrUserIDNotFound
	}

	userID, ok := userIDcontext.(int)
	if !ok {
		return 0, ErrInvalidUserID
	}

	return userID, nil
}
