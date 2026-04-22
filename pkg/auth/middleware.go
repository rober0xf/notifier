package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/domain/entity"
)

func AuthMiddleware(tokenGen TokenGenerator, cookieName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := extractToken(ctx, cookieName)
		if err != nil {
			slog.ErrorContext(ctx.Request.Context(), "failed to extract token", "error", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// validate token
		claims, err := tokenGen.Validate(tokenString)
		if err != nil {
			slog.ErrorContext(ctx.Request.Context(), "failed to verify token", "error", err)
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
		role, err := GetRoleFromContext(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": ErrRoleMissing.Error()})
			return
		}

		if role != entity.RoleAdmin {
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
