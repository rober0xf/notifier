package authentication

import (
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// validatejwt handles cookie and header token
		userID, err := ValidateJWT(c, jwtKey)
		if err != nil {
			// validatejwt already handles the response and abort for all errors
			return
		}

		// set the user id
		c.Set("user_id", userID)

		// if both the token and authentication is valid, the middleware calls the next handler
		c.Next()
	}
}
