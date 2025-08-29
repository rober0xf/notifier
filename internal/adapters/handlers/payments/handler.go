package payments

import (
	"github.com/rober0xf/notifier/internal/services/auth"
	"github.com/rober0xf/notifier/internal/services/payments"
)

type paymentHandler struct {
	PaymentService payments.Service
	Utils          auth.Service
}

func NewPaymentHandler(service payments.Service) *paymentHandler {
	return &paymentHandler{
		PaymentService: service,
	}
}
