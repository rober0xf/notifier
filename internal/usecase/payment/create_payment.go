package payment

import (
	"context"
	"errors"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type CreatePaymentUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewCreatePaymentUseCase(paymentRepo repository.PaymentRepository) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{
		paymentRepo: paymentRepo,
	}
}

func (uc *CreatePaymentUseCase) Execute(ctx context.Context, payment *entity.Payment) (*entity.Payment, error) {
	if err := ValidatePayment(payment); err != nil {
		return nil, err
	}

	if err := uc.paymentRepo.CreatePayment(ctx, payment); err != nil {
		if errors.Is(err, repoErr.ErrAlreadyExists) {
			return nil, domainErr.ErrPaymentAlreadyExists
		}

		return nil, err
	}

	return payment, nil
}
