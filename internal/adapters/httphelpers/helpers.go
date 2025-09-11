package httphelpers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// HTTP handler helper functions
func IsJSONRequest(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}

func IDParameterNotProvided(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter required"})
}

func EmailParameterNotProvided(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "email query parameter required"})
}

func InvalidIDParameter(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "invalid user id: " + err.Error(),
	})
}

// func RequireQuery(query string, f func(*gin.Context, string)) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		value := ctx.Query(query)
// 		if value == "" {
// 			switch query {
// 			case "id":
// 				IDParameterNotProvided(ctx)
// 			case "email":
// 				EmailParameterNotProvided(ctx)
// 			}
// 			return
// 		}
// 		f(ctx, value)
// 	}
// }
