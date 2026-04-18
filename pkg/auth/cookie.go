package auth

import "github.com/gin-gonic/gin"

func SetAuthCookie(c *gin.Context, token string, config CookieConfig) {
	c.SetCookie(
		config.Name,
		token,
		config.MaxAgeSeconds,
		"/",
		"", // empty for current domain
		config.Secure,
		config.HttpOnly,
	)
}

func GetAuthCookie(c *gin.Context, cookieName string) (string, error) {
	token, err := c.Cookie(SessionCookieName)
	if err != nil {
		return "", ErrNoCookie
	}

	return token, nil
}
