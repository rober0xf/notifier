package ports

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

type AuthService interface {
	ParseLoginRequest(c *gin.Context) (dto.LoginRequest, error)
	GetUserIDFromRequest(r *http.Request) (uint, error)
	GenerateToken(userID uint, email string) (string, error)
	ExistsUser(c *gin.Context, email string) (*domain.User, error)
	ValidateToken(tokenString string) (uint, error)
}

type AuthRepository interface {
	ExistsUser(email string) (*domain.User, error)
}
