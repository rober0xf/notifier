package payments

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"gorm.io/gorm"
)

func (p *Payments) Delete(id uint) error {
	var db_payment domain.Payment

	if err := p.db.First(&db_payment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrPaymentNotFound
		}
		return err
	}

	if err := p.db.Delete(&db_payment).Error; err != nil {
		return err
	}

	return nil
}
