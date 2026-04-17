package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/pkg/auth"
)

func (h *UserHandler) GoogleLogin(c *gin.Context) {
	var req struct {
		IDToken string `json:"id_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_token is required"})
		return
	}

	// validate google token
	googleData, err := h.googleVerifier.Verify(req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid google token"})
		return
	}

	out, err := h.oauthUC.Execute(c.Request.Context(), googleData.Sub, googleData.Email, googleData.Name)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrInvalidGoogleID), errors.Is(err, domainErr.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user data"})
		case errors.Is(err, auth.ErrGoogleAccountAlreadyLinked):
			c.JSON(http.StatusConflict, gin.H{"error": "google account already linked"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to login with oauth", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	auth.SetAuthCookie(c, out.Token, auth.CookieConfig{
		Name:            auth.SessionCookieName,
		ExpirationHours: 160,
		Secure:          false,
		HttpOnly:        true,
	})

	c.JSON(http.StatusOK, dto.LoginResponse{
		ID:    out.User.ID,
		Email: out.User.Email,
	})
}
