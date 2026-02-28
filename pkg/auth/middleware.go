package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenGen TokenGenerator, cookieName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := extractToken(ctx, cookieName)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		// validate token
		claims, err := tokenGen.Validate(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Next()
	}
}

// extract the token form the header or cookie
func extractToken(c *gin.Context, cookieName string) (string, error) {
	// first try from the cookie
	if token, err := c.Cookie(cookieName); err == nil {
		return token, nil
	}

	// try from the header
	authHeader := c.GetHeader(AuthHeaderName)
	if authHeader == "" {
		return "", ErrNoToken
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return "", ErrMalformedHeader
	}

	// remove prefix
	return strings.TrimPrefix(authHeader, BearerPrefix), nil
}
