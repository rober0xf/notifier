package payment

import (
	"context"
	"errors"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type UserIDGetter interface {
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
}

type GetAllPaymentsFromUserUseCase struct {
	paymentRepo repository.PaymentRepository
	userRepo    UserIDGetter
}

func NewGetAllPaymentsFromUserUseCase(paymentRepo repository.PaymentRepository, userRepo UserIDGetter) *GetAllPaymentsFromUserUseCase {
	return &GetAllPaymentsFromUserUseCase{
		paymentRepo: paymentRepo,
		userRepo:    userRepo,
	}
}

func (uc *GetAllPaymentsFromUserUseCase) Execute(ctx context.Context, userID int) ([]entity.Payment, error) {
	if userID <= 0 {
		return nil, domainErr.ErrInvalidUserData
	}

	foundUser, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrUserNotFound
		}
		return nil, err
	}

	payments, err := uc.paymentRepo.GetAllPaymentsFromUser(ctx, foundUser.ID)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return []entity.Payment{}, nil // empty list
		}
		return nil, err
	}

	return payments, nil
}
