package payments

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"gorm.io/gorm"
)

func (p *Payments) Get(id uint) (*domain.Payment, error) {
	var payment domain.Payment

	err := p.db.Where("id = ?", id).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, err
	}
	return &payment, nil
}

func (p *Payments) GetAllPayments(user_id uint) ([]*domain.Payment, error) {
	var payments []*domain.Payment

	err := p.db.Where("user_id = ?", user_id).Find(&payments).Error
	if err != nil {
		return nil, err
	}
	if len(payments) == 0 {
		return nil, dto.ErrPaymentNotFound
	}

	return payments, nil
}

func (p *Payments) GetPaymentFromID(id uint, user_id uint) (*domain.Payment, error) {
	var payment domain.Payment

	err := p.db.Where("id = ? AND user_id", id, user_id).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrPaymentNotFound
		}
		return nil, err
	}

	return &payment, nil
}
