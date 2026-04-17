package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	user, err := h.oauthUC.Execute(c.Request.Context(), googleData.Sub, googleData.Email, googleData.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "oauth login failed"})
		return
	}

	token, err := h.tokenGen.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating token"})
		return
	}

	auth.SetAuthCookie(c, token, auth.CookieConfig{
		Name:            auth.SessionCookieName,
		ExpirationHours: 160,
		Secure:          false,
		HttpOnly:        true,
	})

	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"token": token,
	})
}
