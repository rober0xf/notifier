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

type UserIDGetter interface {
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
}

type GetMyPaymentsUseCase struct {
	paymentRepo repository.PaymentRepository
	userRepo    UserIDGetter
}

func NewGetMyPaymentsUseCase(paymentRepo repository.PaymentRepository, userRepo UserIDGetter) *GetMyPaymentsUseCase {
	return &GetMyPaymentsUseCase{
		paymentRepo: paymentRepo,
		userRepo:    userRepo,
	}
}

func (uc *GetMyPaymentsUseCase) Execute(ctx context.Context, userID int) ([]entity.Payment, error) {
	if userID <= 0 {
		return nil, domainErr.ErrInvalidUserData
	}

	foundUser, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrUserNotFound
		}

		return nil, fmt.Errorf("GetMyPaymentsUC.Execute failed to get user %d: %w", userID, err)
	}

	if !foundUser.IsActive {
		return nil, domainErr.ErrUserNotVerified
	}

	payments, err := uc.paymentRepo.GetMyPayments(ctx, foundUser.ID)
	if err != nil {
		return nil, fmt.Errorf("GetMyPaymentsUC.Execute user %d: %w", userID, err)
	}

	return payments, nil
}
