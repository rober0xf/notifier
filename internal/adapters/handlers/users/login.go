package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *userHandler) Login(c *gin.Context) {
	credentials, err := h.Utils.ParseLoginRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Utils.ExistsUser(c, credentials.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if !authentication.VerifyPassword(credentials.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
		return
	}

	token, err := h.Utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while token generation"})
		return
	}

	authentication.SetAuthCookie(c, token)

	c.JSON(http.StatusOK, gin.H{"token": token,
		"user": dto.UserInfo{
			ID:    user.ID,
			Email: user.Email,
		}})
}
