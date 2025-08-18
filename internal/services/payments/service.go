package payments

import (
	"github.com/rober0xf/notifier/internal/ports/payments"
	"gorm.io/gorm"
)

type Payments struct {
	db *gorm.DB
}

func NewPayments(db *gorm.DB) *Payments {
	return &Payments{
		db: db,
	}
}

var _ payments.PaymentService = (*Payments)(nil)
