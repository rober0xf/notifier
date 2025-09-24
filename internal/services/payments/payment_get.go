package payments

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Get(id int) (*domain.Payment, error) {
	return s.Repo.GetPaymentByID(id)
}

func (s *Service) GetAllPayments() ([]domain.Payment, error) {
	return s.Repo.GetAllPayments()
}

func (s *Service) GetAllPaymentsFromUser(email string) ([]domain.Payment, error) {
	if !validateEmail(email) {
		return nil, dto.ErrInvalidPaymentData
	}
	payments, err := s.Repo.GetAllPaymentsFromUser(email)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrNotFound):
			return nil, dto.ErrUserNotFound
		case errors.Is(err, dto.ErrRepository):
			return nil, dto.ErrInternalServerError
		default:
			return nil, dto.ErrInternalServerError
		}
	}
	return payments, nil
}

func validateEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}
