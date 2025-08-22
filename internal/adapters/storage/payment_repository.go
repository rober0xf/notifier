package storage

import (
	"errors"

	"github.com/rober0xf/notifier/internal/domain"
	"github.com/rober0xf/notifier/internal/domain/domain_errors"
	"github.com/rober0xf/notifier/internal/ports"
	"gorm.io/gorm"
)

var _ ports.PaymentRepository = (*Repository)(nil)

func (r *Repository) CreatePayment(payment *domain.Payment) error {
	if err := r.db.Create(payment).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain_errors.ErrNotFound
		}
		return domain_errors.ErrRepository
	}
	return nil
}

func (r *Repository) GetPaymentByID(id uint) (*domain.Payment, error) {
	var payment domain.Payment

	if err := r.db.Where("id = ?", id).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_errors.ErrNotFound
		}
		return nil, domain_errors.ErrRepository
	}
	return &payment, nil
}

func (r *Repository) GetAllPaymentsByUserID(user_id uint) ([]domain.Payment, error) {
	var payments []domain.Payment

	if err := r.db.Where("user_id = ?", user_id).Find(&payments).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_errors.ErrNotFound
		}
		return nil, domain_errors.ErrRepository
	}
	return payments, nil
}

func (r *Repository) GetPaymentByIDAndUserID(id uint, user_id uint) (*domain.Payment, error) {
	var payment domain.Payment

	if err := r.db.Where("id = ? AND user_id = ?", id, user_id).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_errors.ErrNotFound
		}
		return nil, domain_errors.ErrRepository
	}
	return &payment, nil
}

func (r *Repository) UpdatePayment(payment *domain.Payment) error {
	if err := r.db.Save(payment).Error; err != nil {
		return domain_errors.ErrRepository
	}
	return nil
}

func (r *Repository) DeletePayment(id uint) error {
	if err := r.db.Delete(&domain.Payment{}, id).Error; err != nil {
		return domain_errors.ErrRepository
	}
	return nil
}
