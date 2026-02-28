package http

import "github.com/rober0xf/notifier/internal/usecase/payment"

type PaymentHandler struct {
	createPaymentUC          *payment.CreatePaymentUseCase
	getPaymentByIDUC         *payment.GetPaymentByIDUseCase
	getAllPaymentsFromUserUC *payment.GetAllPaymentsFromUserUseCase
	getAllPaymentsUC         *payment.GetAllPaymentsUseCase
	updatePaymentUC          *payment.UpdatePaymentUseCase
	deletePaymentUC          *payment.DeletePaymentUseCase
}

func NewPaymentHandler(
	createPaymentUC *payment.CreatePaymentUseCase,
	getPaymentByIDUC *payment.GetPaymentByIDUseCase,
	getAllPaymentsFromUserUC *payment.GetAllPaymentsFromUserUseCase,
	getAllPaymentsUC *payment.GetAllPaymentsUseCase,
	updatePaymentUC *payment.UpdatePaymentUseCase,
	deletePaymentUC *payment.DeletePaymentUseCase,
) *PaymentHandler {
	return &PaymentHandler{
		createPaymentUC:          createPaymentUC,
		getPaymentByIDUC:         getPaymentByIDUC,
		getAllPaymentsFromUserUC: getAllPaymentsFromUserUC,
		getAllPaymentsUC:         getAllPaymentsUC,
		updatePaymentUC:          updatePaymentUC,
		deletePaymentUC:          deletePaymentUC,
	}
}
