package users

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *userHandler) Login(c *gin.Context) {
	credentials, err := h.Utils.ParseLoginRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	user, err := h.UserService.Repo.GetUserByEmail(c, credentials.Email)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrUserNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		case errors.Is(err, dto.ErrRepository):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

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

	c.SetSameSite(http.SameSiteLaxMode)
	authentication.SetAuthCookie(c, token)

	c.JSON(http.StatusOK, gin.H{"user": dto.LoginResponse{
		Token: token,
		ID:    user.ID,
		Email: user.Email,
	}})
}
