package users

import (
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) GetVerificationEmail(email string) (*domain.User, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// ensure we activate
	user.Active = true
	if err := s.Repo.SetActive(user); err != nil {
		return nil, dto.ErrActivating
	}
	return user, nil
}
