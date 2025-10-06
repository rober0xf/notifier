package payments

import (
	"github.com/rober0xf/notifier/internal/services/auth"
	"github.com/rober0xf/notifier/internal/services/payments"
)

type paymentHandler struct {
	PaymentService *payments.Service
	Utils          *auth.Service
}

func NewPaymentHandler(paymentService *payments.Service, authService *auth.Service) *paymentHandler {
	return &paymentHandler{
		PaymentService: paymentService,
		Utils:          authService,
	}
}
