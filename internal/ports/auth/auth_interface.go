package auth

import (
	"context"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

type AuthService interface {
	ValidateToken(token_string string) (uint, error)
	ParseUserFromToken(token_string string) (*domain.User, error)
	GenerateToken(userID uint, email string) (string, error)
	ExistsUser(ctx context.Context, credentials dto.LoginRequest) (*domain.User, error)
}
