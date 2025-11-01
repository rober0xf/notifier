package payments

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) Get(ctx context.Context, id int) (*domain.Payment, error) {
	payment, err := s.Repo.GetPaymentByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrNotFound):
			return nil, dto.ErrPaymentNotFound
		default:
			return nil, dto.ErrInternalServerError
		}
	}

	return payment, nil
}

func (s *Service) GetAllPayments(ctx context.Context) ([]domain.Payment, error) {
	payments, err := s.Repo.GetAllPayments(ctx)
	if err != nil {
		return nil, dto.ErrInternalServerError
	}

	return payments, nil
}

func (s *Service) GetAllPaymentsFromUser(ctx context.Context, email string) ([]domain.Payment, error) {
	if !validateEmail(email) {
		return nil, dto.ErrInvalidPaymentData
	}
	payments, err := s.Repo.GetAllPaymentsFromUser(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrNotFound):
			return nil, dto.ErrUserNotFound
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
