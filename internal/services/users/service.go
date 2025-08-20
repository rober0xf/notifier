package users

import (
	"github.com/rober0xf/notifier/internal/ports"
)

type Service struct {
	Repo   ports.UserRepository
	jwtKey []byte
}

func NewUsers(repo ports.UserRepository, jwtKey []byte) *Service {
	return &Service{
		Repo:   repo,
		jwtKey: jwtKey,
	}
}

var _ ports.UserService = (*Service)(nil)
