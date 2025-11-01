package ports

import (
	"context"

	"github.com/rober0xf/notifier/internal/domain"
)

type PaymentService interface {
	Create(ctx context.Context, user *domain.Payment) (*domain.Payment, error)
	Get(ctx context.Context, id int) (*domain.Payment, error)
	GetAllPayments(ctx context.Context) ([]domain.Payment, error)
	GetAllPaymentsFromUser(ctx context.Context, email string) ([]domain.Payment, error)
	Update(ctx context.Context, id int, payment *domain.UpdatePayment) (*domain.Payment, error)
	Delete(ctx context.Context, id int) error
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	GetAllPayments(ctx context.Context) ([]domain.Payment, error)
	GetPaymentByID(ctx context.Context, id int) (*domain.Payment, error)
	GetAllPaymentsFromUser(ctx context.Context, email string) ([]domain.Payment, error)
	UpdatePayment(ctx context.Context, payment *domain.Payment) error
	DeletePayment(ctx context.Context, id int) error
}
