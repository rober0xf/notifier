package payments

import (
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Update(payment *domain.Payment) (*domain.Payment, error) {
	if err := s.Repo.UpdatePayment(payment); err != nil {
		return nil, dto.ErrInternalServerError
	}
	return payment, nil
}
