package services

import (
	"errors"
	"strings"

	"github.com/rober0xf/notifier/internal/core/shared"
	"github.com/rober0xf/notifier/internal/models"
	"gorm.io/gorm"
)

type PaymentService struct {
	db *gorm.DB
}

func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{
		db: db,
	}
}

func (ps *PaymentService) CreatePaymentService(payment *models.Payment) error {
	if payment == nil {
		return errors.New("payment is nil")
	}

	if payment.UserID == 0 || payment.NetAmount <= 0 || payment.GrossAmount <= 0 || payment.Name == "" || payment.Type == "" || payment.Date.IsZero() {
		return shared.ErrInvalidPaymentData
	}

	input_payment := models.Payment{
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

	_, err := ps.RegisterPayment(&input_payment)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PaymentService) GetPaymentService(id uint) (*models.Payment, error) {
	var payment models.Payment

	err := ps.db.Where("id = ?", id).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrPaymentNotFound
		}
		return nil, err
	}

	return &payment, nil
}

func (ps *PaymentService) GetAllPaymentsService(user_id uint) ([]*models.Payment, error) {
	var payments []*models.Payment

	err := ps.db.Where("user_id = ?", user_id).Find(&payments).Error
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (ps *PaymentService) GetPaymentFromIDService(id uint, user_id uint) (*models.Payment, error) {
	var payment models.Payment

	err := ps.db.Where("id = ? AND user_id", id, user_id).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrPaymentNotFound
		}
		return nil, err
	}

	return &payment, nil
}

func (ps *PaymentService) UpdatePaymentService(payment *models.Payment) (*models.Payment, error) {
	var db_payment models.Payment
	err := ps.db.Where("id = ? AND user_id = ?", payment.ID, payment.UserID).First(&db_payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrPaymentNotFound
		}
		return nil, err
	}

	if payment.UserID == 0 || payment.NetAmount <= 0 || payment.GrossAmount <= 0 || payment.Name == "" || payment.Type == "" || payment.Date.IsZero() {
		return nil, shared.ErrInvalidPaymentData
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

	if err := ps.db.Save(&db_payment).Error; err != nil {
		return nil, err
	}

	return &db_payment, nil
}

func (ps *PaymentService) DeletePaymentService(id uint) error {
	var db_payment models.Payment

	if err := ps.db.First(&db_payment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shared.ErrPaymentNotFound
		}
		return err
	}

	if err := ps.db.Delete(&db_payment).Error; err != nil {
		return err
	}

	return nil
}

// ------------- helper functions ------------------------
func (ps *PaymentService) RegisterPayment(payment *models.Payment) (*models.Payment, error) {
	if err := ps.db.Create(payment).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "null") || strings.Contains(err.Error(), "invalid") {
			return nil, shared.ErrInvalidPaymentData
		}
		return nil, err
	}

	return payment, nil
}
