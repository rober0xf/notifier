package ports

import "github.com/rober0xf/notifier/internal/domain"

type PaymentService interface {
	Create(*domain.Payment) (*domain.Payment, error)
	Get(id uint) (*domain.Payment, error)
	GetAllPayments() ([]domain.Payment, error)
	GetAllPaymentsFromUser(email string) ([]domain.Payment, error)
	Update(*domain.Payment) (*domain.Payment, error)
	Delete(id uint) error
}

type PaymentRepository interface {
	CreatePayment(payment *domain.Payment) error
	GetAllPayments() ([]domain.Payment, error)
	GetPaymentByID(id uint) (*domain.Payment, error)
	GetAllPaymentsFromUser(email string) ([]domain.Payment, error)
	UpdatePayment(payment *domain.Payment) error
	DeletePayment(id uint) error
}
