package payments

import (
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (s *Service) Delete(id uint) error {
	if err := s.Repo.DeletePayment(id); err != nil {
		return dto.ErrInternalServerError
	}
	return nil
}
