package auth

import (
	"github.com/rober0xf/notifier/internal/ports"
)

type Service struct {
	UserRepo ports.UserRepository
	jwtKey   []byte
}

func NewAuthService(repo ports.UserRepository, jwtKey []byte) *Service {
	return &Service{
		UserRepo: repo,
		jwtKey:   jwtKey,
	}
}

var _ ports.AuthService = (*Service)(nil)
