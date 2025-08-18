package payments

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"gorm.io/gorm"
)

func (p *Payments) Update(payment *domain.Payment) (*domain.Payment, error) {
	var db_payment domain.Payment
	err := p.db.Where("id = ? AND user_id = ?", payment.ID, payment.UserID).First(&db_payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrPaymentNotFound
		}
		return nil, err
	}

	if payment.UserID == 0 || payment.NetAmount <= 0 || payment.GrossAmount <= 0 || payment.Name == "" || payment.Type == "" || payment.Date.IsZero() {
		return nil, dto.ErrInvalidPaymentData
	}

	// we need to update it manually
	db_payment.NetAmount = payment.NetAmount
	db_payment.GrossAmount = payment.GrossAmount
	db_payment.Deductible = payment.Deductible
	db_payment.Name = payment.Name
	db_payment.Type = payment.Type
	db_payment.Date = payment.Date
	db_payment.Recurrent = payment.Recurrent
	db_payment.Paid = payment.Paid

	if err := p.db.Save(&db_payment).Error; err != nil {
		return nil, err
	}

	return &db_payment, nil
}
