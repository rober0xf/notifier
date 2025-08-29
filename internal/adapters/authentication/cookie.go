package authentication

import (
	"github.com/gin-gonic/gin"
)

const (
	TokenExpirationHours = 6
	BearerPrefix         = "BEARER "
	SessionCookieName    = "session_token"
	AuthHeaderName       = "Authorization"
)

func SetAuthCookie(c *gin.Context, token string) {
	c.SetCookie(SessionCookieName,
		token,
		int(TokenExpirationHours),
		"/",
		"", // empty for current domain
		true,
		true,
	)
}
