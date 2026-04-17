package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	authErr "github.com/rober0xf/notifier/pkg/auth"
)

func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Param("token")

	user, err := h.verifyEmailUC.Execute(c.Request.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, authErr.ErrInvalidToken):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid or expired verification link"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to verify user", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully", "user": toUserResponse(user)})
}
