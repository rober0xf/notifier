package storage

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
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

func (r *Repository) GetAllPayments() ([]domain.Payment, error) {
	var payments []domain.Payment

	if err := r.db.Find(&payments).Error; err != nil {
		return nil, domain_errors.ErrRepository
	}
	return payments, nil
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

func (r *Repository) GetAllPaymentsFromUser(email string) ([]domain.Payment, error) {
	// first we need to check if the email exists in our db
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, domain_errors.ErrRepository
	}

	var payments []domain.Payment
	if err := r.db.Where("user_id = ?", user.ID).Find(&payments).Error; err != nil {
		return nil, domain_errors.ErrRepository
	}

	// if the user doesnt have any payments it returns []
	return payments, nil
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
