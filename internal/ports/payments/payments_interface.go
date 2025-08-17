package payments

import "github.com/rober0xf/notifier/internal/domain"

type PaymentService interface {
	Create(*domain.Payment) error
	Get(id uint) (*domain.Payment, error)
	GetAllPayments(user_id uint) ([]*domain.Payment, error)
	GetPaymentFromID(id uint, user_id uint) (*domain.Payment, error)
	Update(*domain.Payment) (*domain.Payment, error)
	Delete(id uint) error
}
