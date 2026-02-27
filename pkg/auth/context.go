package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrUserIDNotFound = errors.New("user_id not found in context")
var ErrInvalidUserID = errors.New("invalid user_id type")

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
