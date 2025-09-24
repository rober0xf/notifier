package ports

import "github.com/rober0xf/notifier/internal/domain"

type PaymentService interface {
	Create(*domain.Payment) (*domain.Payment, error)
	Get(id int) (*domain.Payment, error)
	GetAllPayments() ([]domain.Payment, error)
	GetAllPaymentsFromUser(email string) ([]domain.Payment, error)
	Update(id int, payment *domain.UpdatePayment) (*domain.Payment, error)
	Delete(id int) error
}

type PaymentRepository interface {
	CreatePayment(payment *domain.Payment) error
	GetAllPayments() ([]domain.Payment, error)
	GetPaymentByID(id int) (*domain.Payment, error)
	GetAllPaymentsFromUser(email string) ([]domain.Payment, error)
	UpdatePayment(payment *domain.Payment) error
	DeletePayment(id int) error
}
