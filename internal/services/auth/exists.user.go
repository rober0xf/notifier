package auth

import (
	"context"

	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) ExistsUser(ctx context.Context, credentials dto.LoginRequest) (*domain.User, error) {
	user, err := s.Repo.ExistsUser(ctx, credentials)
	if err != nil {
		return nil, err
	}

	hashed_password, err := authentication.HashPassword(credentials.Password)
	if err != nil {
		return nil, dto.ErrInternalServerError
	}

	if user.Password != hashed_password {
		return nil, dto.ErrInvalidCredentials
	}

	return user, nil
}
