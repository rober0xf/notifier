package ports

import (
	"context"
	"net/http"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

type AuthService interface {
	ParseLoginRequest(r *http.Request) (dto.LoginRequest, error)
	GetUserIDFromRequest(r *http.Request) (uint, error)
	GenerateToken(userID uint, email string) (string, error)
	ExistsUser(ctx context.Context, credentials dto.LoginRequest) (*domain.User, error)
}

type AuthRepository interface {
	ValidateToken(token_string string) (uint, error)
	ParseUserFromToken(token_string string) (*domain.User, error)
	GenerateToken(userID uint, email string) (string, error)
	ExistsUser(ctx context.Context, credentials dto.LoginRequest) (*domain.User, error)
}
