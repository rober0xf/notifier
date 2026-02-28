package payment

import (
	"context"
	"errors"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type GetPaymentByIDUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewGetPaymentByIDUseCase(paymentRepo repository.PaymentRepository) *GetPaymentByIDUseCase {
	return &GetPaymentByIDUseCase{
		paymentRepo: paymentRepo,
	}
}

func (uc *GetPaymentByIDUseCase) Execute(ctx context.Context, id int) (*entity.Payment, error) {
	if id <= 0 {
		return nil, domainErr.ErrInvalidPaymentData
	}

	payment, err := uc.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrPaymentNotFound
		}
		return nil, err
	}

	return payment, nil
}
