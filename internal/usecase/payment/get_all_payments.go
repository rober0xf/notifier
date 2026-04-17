package payment

import (
	"context"

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
	return uc.paymentRepo.GetAllPayments(ctx)
}
