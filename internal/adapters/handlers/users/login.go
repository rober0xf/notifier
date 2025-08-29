package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *userHandler) Login(c *gin.Context) {
	credentials, err := h.Utils.ParseLoginRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.Utils.ExistsUser(c, credentials)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
		return
	}

	token, err := h.Utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while token generation"})
		return
	}

	authentication.SetAuthCookie(c, token)

	// if it comes from json
	if httphelpers.IsJSONRequest(c.Request) {
		c.JSON(http.StatusOK, gin.H{"token": token,
			"user": dto.UserInfo{
				ID:    user.ID,
				Email: user.Email,
			}})
	} else {
		// from frontend
		c.Header("HX-Redirect", "/dashboard")
		c.Status(http.StatusOK)
	}
}
