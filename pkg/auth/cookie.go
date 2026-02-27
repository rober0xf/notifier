package auth

import "github.com/gin-gonic/gin"

type CookieConfig struct {
	Name            string
	ExpirationHours int
	Secure          bool // true in https
	HttpOnly        bool
	SameSite        int
}

func SetAuthCookie(c *gin.Context, token string, config CookieConfig) {
	c.SetCookie(
		config.Name,
		token,
		config.ExpirationHours,
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
