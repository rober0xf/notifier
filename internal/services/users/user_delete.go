package users

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain/domain_errors"
)

func (s *Service) Delete(id uint) error {

	_, err := s.Repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, domain_errors.ErrNotFound) {
			return dto.ErrUserNotFound
		}
		return dto.ErrInternalServerError
	}

	if err := s.Repo.DeleteUser(id); err != nil {
		return dto.ErrInternalServerError
	}
	return nil
}
