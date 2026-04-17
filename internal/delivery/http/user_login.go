package http

import (
	"errors"
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

	user, err := h.loginUC.Execute(c.Request.Context(), payload.Email, payload.Password)
	if err != nil {
		switch {
		case errors.Is(err, authErr.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		case errors.Is(err, authErr.ErrEmailNotVerified):
			c.JSON(http.StatusForbidden, gin.H{"error": "email not verified"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	token, err := h.tokenGen.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": auth.ErrGeneratingToken})
		return
	}

	auth.SetAuthCookie(c, token, auth.CookieConfig{
		Name:            auth.SessionCookieName,
		ExpirationHours: 160,
		Secure:          false,
		HttpOnly:        true,
	})

	c.JSON(http.StatusOK, dto.LoginResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}
