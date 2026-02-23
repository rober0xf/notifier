package repository

import (
	"context"

	"github.com/rober0xf/notifier/internal/domain/entity"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	GetAllPayments(ctx context.Context) ([]entity.Payment, error)
	GetPaymentByID(ctx context.Context, id int) (*entity.Payment, error)
	GetAllPaymentsFromUser(ctx context.Context, userID int) ([]entity.Payment, error)
	UpdatePayment(ctx context.Context, payment *entity.Payment) error
	DeletePayment(ctx context.Context, id int) error
}

// type PaymentService interface {
// 	Create(ctx context.Context, user *entity.Payment) (*entity.Payment, error)
// 	Get(ctx context.Context, id int) (*entity.Payment, error)
// 	GetAllPayments(ctx context.Context) ([]entity.Payment, error)
// 	GetAllPaymentsFromUser(ctx context.Context, email string) ([]entity.Payment, error)
// 	Update(ctx context.Context, id int, payment *dto.UpdatePayment) (*entity.Payment, error)
// 	Delete(ctx context.Context, id int) error
// }
