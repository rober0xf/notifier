package dto

import (
	"fmt"
	"time"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	uc "github.com/rober0xf/notifier/internal/usecase/payment"
)

type CreatePaymentRequest struct {
	Name       string                 `json:"name" binding:"required,min=3,max=100"`
	Amount     float64                `json:"amount" binding:"required,gt=0"`
	Type       entity.TransactionType `json:"type" binding:"required"`
	Category   entity.CategoryType    `json:"category" binding:"required"`
	Date       string                 `json:"date" binding:"required,datetime=2006-01-02"`
	DueDate    string                 `json:"due_date" binding:"omitempty,datetime=2006-01-02"`
	Paid       bool                   `json:"paid"`
	PaidAt     string                 `json:"paid_at" binding:"omitempty,datetime=2006-01-02"`
	Recurrent  bool                   `json:"recurrent"`
	Frequency  entity.FrequencyType   `json:"frequency" binding:"omitempty"`
	ReceiptURL string                 `json:"receipt_url" binding:"omitempty,url"`
}

type PaymentResponse struct {
	ID         int32                  `json:"id"`
	Name       string                 `json:"name"`
	Amount     float64                `json:"amount"`
	Type       entity.TransactionType `json:"type"`
	Category   entity.CategoryType    `json:"category"`
	Date       string                 `json:"date"`
	DueDate    *string                `json:"due_date,omitempty"`
	Paid       bool                   `json:"paid"`
	PaidAt     *string                `json:"paid_at,omitempty"`
	Recurrent  bool                   `json:"recurrent"`
	Frequency  *entity.FrequencyType  `json:"frequency,omitempty"`
	ReceiptURL *string                `json:"receipt_url,omitempty"`
}

type UpdatePaymentRequest struct {
	Name       *string                 `json:"name,omitempty"`
	Amount     *float64                `json:"amount,omitempty"`
	Type       *entity.TransactionType `json:"type,omitempty"`
	Category   *entity.CategoryType    `json:"category,omitempty"`
	Date       *string                 `json:"date,omitempty"`
	DueDate    *string                 `json:"due_date,omitempty"`
	Paid       *bool                   `json:"paid,omitempty"`
	PaidAt     *string                 `json:"paid_at,omitempty"`
	Recurrent  *bool                   `json:"recurrent,omitempty"`
	Frequency  *entity.FrequencyType   `json:"frequency,omitempty"`
	ReceiptURL *string                 `json:"receipt_url,omitempty"`
}

func (p *CreatePaymentRequest) Validate() error {
	if !p.Category.IsValid() {
		return fmt.Errorf("invalid category: %w", domainErr.ErrInvalidCategory)
	}

	if !p.Type.IsValid() {
		return fmt.Errorf("invalid transaction type: %w", domainErr.ErrInvalidTransactionType)
	}

	if !p.Frequency.IsValid() {
		return fmt.Errorf("invalid frequency type: %w", domainErr.ErrInvalidFrequency)
	}

	if p.Recurrent && p.Frequency == "" {
		return fmt.Errorf("frequency is required for recurrent payments: %w", domainErr.ErrInvalidFrequency)
	}

	if p.Frequency != "" && !p.Recurrent {
		return fmt.Errorf("recurrent is required if frequency is set: %w", domainErr.ErrInvalidFrequency)
	}

	if p.PaidAt != "" && !p.Paid {
		return fmt.Errorf("piad is required if paid_at is set: %w", domainErr.ErrInvalidPaymentData)
	}

	if p.Date != "" && p.DueDate != "" {
		date, _ := time.Parse("2006-01-02", p.Date)
		dueDate, _ := time.Parse("2006-01-02", p.DueDate)
		if dueDate.Before(date) {
			return fmt.Errorf("due_date cannot be before payment date: %w", domainErr.ErrInvalidDate)
		}
	}

	return nil
}

func (p *UpdatePaymentRequest) Validate() error {
	if p.Paid != nil && !*p.Paid {
		if p.PaidAt != nil && *p.PaidAt != "" {
			return fmt.Errorf("cannot set paid_at if payment is not paid: %w", domainErr.ErrInvalidPaymentData)
		}
		p.PaidAt = nil
	}

	if p.Recurrent != nil && *p.Recurrent && p.Frequency != nil && *p.Frequency == "" {
		return fmt.Errorf("frequency is required for recurrent payments: %w", domainErr.ErrInvalidFrequency)
	}

	if p.Date != nil && *p.Date != "" && p.DueDate != nil && *p.DueDate != "" {
		date, err1 := time.Parse("2006-01-02", *p.Date)
		dueDate, err2 := time.Parse("2006-01-02", *p.DueDate)
		if err1 == nil && err2 == nil && dueDate.Before(date) {
			return fmt.Errorf("due_date cannot be before date: %w", domainErr.ErrInvalidDate)
		}
	}

	return nil
}

func (p *UpdatePaymentRequest) ToInput() uc.UpdatePaymentInput {
	return uc.UpdatePaymentInput{
		Name:       p.Name,
		Amount:     p.Amount,
		Type:       p.Type,
		Category:   p.Category,
		Date:       p.Date,
		DueDate:    p.DueDate,
		Paid:       p.Paid,
		PaidAt:     p.PaidAt,
		Recurrent:  p.Recurrent,
		Frequency:  p.Frequency,
		ReceiptURL: p.ReceiptURL,
	}
}
