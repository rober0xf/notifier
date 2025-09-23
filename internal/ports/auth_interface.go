package ports

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

type AuthService interface {
	ParseLoginRequest(c *gin.Context) (dto.LoginRequest, error)
	GetUserIDFromRequest(r *http.Request) (int, error)
	GenerateToken(userID int, email string) (string, error)
	ValidateToken(tokenString string, jwtKey []byte) (int, error)
}
