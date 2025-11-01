package users

import (
	"context"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) GetVerificationEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// ensure we activate
	user.Active = true
	if err := s.Repo.UpdateUserActive(ctx, user.ID, true); err != nil {
		return nil, dto.ErrActivating
	}

	return user, nil
}
