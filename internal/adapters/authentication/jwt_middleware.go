package authentication

import (
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// validatejwt handles cookie and header token
		_, err := ValidateJWT(c, jwtKey)
		if err != nil {
			// validatejwt already handles the response and abort for all errors
			return
		}

		// if both the token and authentication is valid, the middleware calls the next handler
		c.Next()
	}
}
