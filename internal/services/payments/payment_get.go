package payments

import (
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Get(id uint) (*domain.Payment, error) {
	return s.Repo.GetPaymentByID(id)
}

func (s *Service) GetAllPayments() ([]domain.Payment, error) {
	return s.Repo.GetAllPayments()
}

func (s *Service) GetAllPaymentsFromUser(email string) ([]domain.Payment, error) {
	return s.Repo.GetAllPaymentsFromUser(email)
}
