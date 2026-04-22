package payment

import (
	"context"
	"errors"
	"fmt"

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
	if payment.UserID <= 0 {
		return nil, domainErr.ErrInvalidPaymentData
	}

	created, err := uc.paymentRepo.CreatePayment(ctx, payment)
	if err != nil {
		if errors.Is(err, repoErr.ErrAlreadyExists) {
			return nil, domainErr.ErrPaymentAlreadyExists
		}

		return nil, fmt.Errorf("CreatePaymentUC.Execute failed to create payment: %w", err)
	}

	return created, nil
}
