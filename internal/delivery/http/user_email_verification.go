package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	authErr "github.com/rober0xf/notifier/pkg/auth"
)

func (h *UserHandler) GetVerificationEmail(c *gin.Context) {
	email := c.Param("email")
	token := c.Param("hash")

	if email == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification link"})
		return
	}

	user, err := h.verifyEmailUC.Execute(c.Request.Context(), email, token)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, authErr.ErrInvalidToken):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired verification link"})
		case errors.Is(err, domainErr.ErrAlreadyVerified):
			c.JSON(http.StatusConflict, gin.H{"error": "email already verified"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully", "user": toUserResponse(user)})
}
