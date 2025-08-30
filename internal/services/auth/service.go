package auth

import (
	"github.com/rober0xf/notifier/internal/ports"
)

type Service struct {
	AuthRepo ports.AuthRepository
	UserRepo ports.UserRepository
	jwtKey   []byte
}

func NewAuthService(authRepo ports.AuthRepository, userRepo ports.UserRepository, jwtKey []byte) *Service {
	return &Service{
		AuthRepo: authRepo,
		UserRepo: userRepo,
		jwtKey:   jwtKey,
	}
}

var _ ports.AuthService = (*Service)(nil)
