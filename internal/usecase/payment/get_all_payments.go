package payment

import (
	"context"
	"errors"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type GetAllPaymentsUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewGetAllPaymentsUseCase(paymentRepo repository.PaymentRepository) *GetAllPaymentsUseCase {
	return &GetAllPaymentsUseCase{
		paymentRepo: paymentRepo,
	}
}

func (uc *GetAllPaymentsUseCase) Execute(ctx context.Context) ([]entity.Payment, error) {
	payments, err := uc.paymentRepo.GetAllPayments(ctx)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return []entity.Payment{}, nil // empty list
		}
		return nil, err
	}

	return payments, nil
}
