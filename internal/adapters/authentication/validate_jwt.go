package authentication

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

const (
	AuthHeaderName = "Authorization"
)

func ValidateJWT(c *gin.Context, jwtKey []byte) (uint, error) {
	// first we try to get the token from the cookie
	tokenString, err := getAuthCookie(c)
	if err != nil {
		// if the there is no cookie then we get the token from the header
		authHeader := c.GetHeader(AuthHeaderName)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
			c.Abort()
			return 0, dto.ErrNoToken
		}

		// remove the prefix
		if len(authHeader) > 7 && authHeader[:7] == BearerPrefix {
			tokenString = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "malformed authorization header"})
			c.Abort()
			return 0, dto.ErrMalformedHeader
		}
	}

	userID, err := validateTokenString(tokenString, jwtKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		c.Abort()
		return 0, err
	}
	
	return userID, nil
}
