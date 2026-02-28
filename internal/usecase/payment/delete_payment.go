package payment

import (
	"context"
	"errors"

	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type DeletePaymentUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewDeletePaymentUseCase(paymentRepo repository.PaymentRepository) *DeletePaymentUseCase {
	return &DeletePaymentUseCase{
		paymentRepo: paymentRepo,
	}
}

func (uc *DeletePaymentUseCase) Execute(ctx context.Context, id int) error {
	if id <= 0 {
		return domainErr.ErrInvalidPaymentData
	}

	if err := uc.paymentRepo.DeletePayment(ctx, id); err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return domainErr.ErrPaymentNotFound
		}
		return err
	}

	return nil
}
