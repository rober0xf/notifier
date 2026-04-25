package payment

import (
	"context"
	"errors"
	"fmt"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	authErr "github.com/rober0xf/notifier/pkg/auth"
)

type DeletePaymentUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewDeletePaymentUseCase(paymentRepo repository.PaymentRepository) *DeletePaymentUseCase {
	return &DeletePaymentUseCase{
		paymentRepo: paymentRepo,
	}
}

func (uc *DeletePaymentUseCase) Execute(ctx context.Context, id, userID int, userRole entity.Role) error {
	payment, err := uc.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return domainErr.ErrPaymentNotFound
		}

		return fmt.Errorf("failed getting payment by id: %w", err)
	}

	if payment.UserID != userID && userRole != entity.RoleAdmin {
		return authErr.ErrForbidden
	}

	if err := uc.paymentRepo.DeletePayment(ctx, id); err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return domainErr.ErrPaymentNotFound
		}

		return fmt.Errorf("DeletePaymentUC.Execute failed to delete payment: %w", err)
	}

	return nil
}
