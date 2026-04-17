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
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// validate token
		claims, err := tokenGen.Validate(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, ok := ctx.Get("role")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": ErrRoleMissing.Error()})
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": ErrInvalidRole.Error()})
			return
		}

		if roleStr != "admin" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

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
