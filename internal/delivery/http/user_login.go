package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/pkg/auth"
	authErr "github.com/rober0xf/notifier/pkg/auth"
)

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	// just in case
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
	}

	token, user, err := h.loginUC.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, authErr.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"}) // safer
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	auth.SetAuthCookie(c, token, auth.CookieConfig{
		Name:            auth.SessionCookieName,
		ExpirationHours: 160,
		Secure:          false,
		HttpOnly:        true,
	})

	c.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
		ID:    user.ID,
		Email: user.Email,
	})
}
