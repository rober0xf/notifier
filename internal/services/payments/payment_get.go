package payments

import (
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Get(id uint) (*domain.Payment, error) {
	payment, err := s.Repo.GetPaymentByID(id)
	if err != nil {
		return nil, dto.ErrInternalServerError
	}
	return payment, nil
}

func (s *Service) GetAllPayments(user_id uint) ([]domain.Payment, error) {
	payments, err := s.Repo.GetAllPaymentsByUserID(user_id)
	if err != nil {
		return nil, dto.ErrInternalServerError
	}
	if len(payments) == 0 {
		return nil, dto.ErrPaymentNotFound
	}
	return payments, nil
}

func (s *Service) GetPaymentFromIDAndUserID(id uint, user_id uint) (*domain.Payment, error) {
	payment, err := s.Repo.GetPaymentByIDAndUserID(id, user_id)
	if err != nil {
		return nil, dto.ErrInternalServerError
	}
	return payment, nil
}
