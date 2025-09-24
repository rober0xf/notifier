package payments

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (s *Service) Delete(id int) error {
	if err := s.Repo.DeletePayment(id); err != nil {
		switch {
		case errors.Is(err, dto.ErrNotFound):
			return dto.ErrPaymentNotFound
		case errors.Is(err, dto.ErrRepository):
			return dto.ErrInternalServerError
		default:
			return dto.ErrInternalServerError
		}
	}
	return nil
}
