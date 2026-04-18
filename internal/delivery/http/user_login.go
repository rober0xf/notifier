package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/pkg/auth"
	authErr "github.com/rober0xf/notifier/pkg/auth"
)

func (h *UserHandler) Login(c *gin.Context) {
	var payload dto.LoginPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	out, err := h.loginUC.Execute(c.Request.Context(), payload.Email, payload.Password)
	if err != nil {
		switch {
		case errors.Is(err, authErr.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		case errors.Is(err, authErr.ErrEmailNotVerified):
			c.JSON(http.StatusForbidden, gin.H{"error": "email not verified"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to login user", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	auth.SetAuthCookie(c, out.Token, auth.CookieConfig{
		Name:          auth.SessionCookieName,
		MaxAgeSeconds: 24 * 3600,
		Secure:        false,
		HttpOnly:      true,
	})

	c.JSON(http.StatusOK, dto.LoginResponse{
		ID:    out.User.ID,
		Email: out.User.Email,
	})
}
