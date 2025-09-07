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
