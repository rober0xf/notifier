package auth

import "github.com/rober0xf/notifier/internal/ports"

type Service struct {
	Repo   ports.AuthRepository
	jwtKey []byte
}

func NewAuthService(repo ports.AuthRepository, jwtKey []byte) *Service {
	return &Service{
		Repo:   repo,
		jwtKey: jwtKey,
	}
}

var _ ports.AuthService = (*Service)(nil)
