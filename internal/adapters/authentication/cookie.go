package authentication

import (
	"github.com/gin-gonic/gin"
)

const (
	TokenExpirationHours = 6
	BearerPrefix         = "Bearer "
	SessionCookieName    = "session_token"
)

func SetAuthCookie(c *gin.Context, token string) {
	c.SetCookie(
		SessionCookieName,
		token,
		int(TokenExpirationHours*3600),
		"",
		"", // empty for current domain
		false,
		true,
	)
}

func getAuthCookie(c *gin.Context) (string, error) {
	tokenString, err := c.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
