package payment

import (
	"context"
	"fmt"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/domain/repository"
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
		return nil, fmt.Errorf("GetAllPaymentsUC.Execute failed to get all payments: %w", err)
	}

	return payments, nil
}
