package api

import (
	routes "github.com/rober0xf/notifier/internal/delivery/http"
	"github.com/rober0xf/notifier/internal/domain/repository"
	"github.com/rober0xf/notifier/internal/usecase/payment"
)

func buildPaymentRoutes(paymentRepo repository.PaymentRepository, userRepo payment.UserIDGetter) *routes.PaymentHandler {
	return routes.NewPaymentHandler(
		payment.NewCreatePaymentUseCase(paymentRepo),
		payment.NewGetPaymentByIDUseCase(paymentRepo),
		payment.NewGetAllPaymentsFromUserUseCase(paymentRepo, userRepo),
		payment.NewGetAllPaymentsUseCase(paymentRepo),
		payment.NewUpdatePaymentUseCase(paymentRepo),
		payment.NewDeletePaymentUseCase(paymentRepo),
	)
}
