package users

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	domainErrors "github.com/rober0xf/notifier/internal/domain/errors"
)

func (s *Service) Delete(id uint) error {

	_, err := s.Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, domainErrors.ErrNotFound) {
			return dto.ErrUserNotFound
		}
		return dto.ErrInternalServerError
	}

	if err := s.Repo.Delete(id); err != nil {
		return dto.ErrInternalServerError
	}
	return nil
}
