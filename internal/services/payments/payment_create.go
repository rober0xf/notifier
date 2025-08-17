package payments

import (
	"errors"
	"strings"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (p *Payments) Create(payment *domain.Payment) error {
	if payment == nil {
		return errors.New("payment is nil")
	}

	if payment.UserID == 0 || payment.NetAmount <= 0 || payment.GrossAmount <= 0 || payment.Name == "" || payment.Type == "" || payment.Date.IsZero() {
		return dto.ErrInvalidPaymentData
	}

	input_payment := domain.Payment{
		UserID:      payment.UserID,
		NetAmount:   payment.NetAmount,
		GrossAmount: payment.GrossAmount,
		Deductible:  payment.Deductible,
		Name:        payment.Name,
		Type:        payment.Type,
		Date:        payment.Date,
		Recurrent:   payment.Recurrent,
		Paid:        payment.Paid,
	}

	_, err := p.Register(&input_payment)
	if err != nil {
		return err
	}

	return nil
}

func (p *Payments) Register(payment *domain.Payment) (*domain.Payment, error) {
	if err := p.db.Create(payment).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "null") || strings.Contains(err.Error(), "invalid") {
			return nil, dto.ErrInvalidPaymentData
		}
		return nil, err
	}

	return payment, nil
}
