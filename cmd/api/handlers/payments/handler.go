package payments

import (
	"github.com/go-playground/validator/v10"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/ports/payments"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Handler struct {
	PaymentService payments.PaymentService
	AuthUtils      httphelpers.AuthHelper
}

func NewPaymentHandler(paymentService payments.PaymentService, authUtils httphelpers.AuthHelper) *Handler {
	return &Handler{
		PaymentService: paymentService,
		AuthUtils:      authUtils,
	}
}

type input_payment struct {
	NetAmount   float64 `json:"net_amount" validate:"required"`
	GrossAmount float64 `json:"gross_amount"`
	Deductible  float64 `json:"deductible"`
	Name        string  `gorm:"not null" json:"name" validate:"required"`
	Type        string  `gorm:"not null" json:"type" validate:"required"`
	Date        string  `json:"date" validate:"required"`
	Recurrent   bool    `gorm:"not null" json:"recurrent" validate:"required"`
	Paid        bool    `gorm:"not null" json:"paid" validate:"required"`
}
