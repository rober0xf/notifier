package dto

import (
	"fmt"
	"time"

	"github.com/rober0xf/notifier/internal/domain/entity"
)

type UpdatePaymentRequest struct {
	Name       *string                 `json:"name,omitempty"`
	Amount     *float64                `json:"amount,omitempty"`
	Type       *entity.TransactionType `json:"type,omitempty"`
	Category   *entity.CategoryType    `json:"category,omitempty"`
	Date       *string                 `json:"date,omitempty"`
	DueDate    *string                 `json:"due_date,omitempty"`
	Paid       *bool                   `json:"paid"`
	PaidAt     *string                 `json:"paid_at,omitempty"`
	Recurrent  *bool                   `json:"recurrent,omitempty"`
	Frequency  *entity.FrequencyType   `json:"frequency,omitempty"`
	ReceiptURL *string                 `json:"receipt_url,omitempty"`
}

type CreatePaymentRequest struct {
	Name       string                 `json:"name" binding:"required,min=3,max=100"`
	Amount     float64                `json:"amount" binding:"required,gt=0"`
	Type       entity.TransactionType `json:"type" binding:"required,oneof=expense income subscription"`
	Category   entity.CategoryType    `json:"category" binding:"required,oneof=electronics entertainment education clothing work sports"`
	Date       string                 `json:"date" binding:"required,datetime=2006-01-02"`
	DueDate    string                 `json:"due_date" binding:"omitempty,datetime=2006-01-02"`
	Paid       bool                   `json:"paid"`
	PaidAt     string                 `json:"paid_at" binding:"omitempty,datetime=2006-01-02"`
	Recurrent  bool                   `json:"recurrent"`
	Frequency  entity.FrequencyType   `json:"frequency" binding:"omitempty,oneof=daily weekly monthly yearly"`
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

func (p *CreatePaymentRequest) Validate() error {
	if p.Recurrent && p.Frequency == "" {
		return fmt.Errorf("frequency is required for recurrent payments")
	}

	if p.Date != "" && p.DueDate != "" {
		date, _ := time.Parse("2006-01-02", p.Date)
		dueDate, _ := time.Parse("2006-01-02", p.DueDate)
		if dueDate.Before(date) {
			return fmt.Errorf("due_date cannot be before payment date")
		}
	}

	return nil
}

func (p *UpdatePaymentRequest) Validate() error {
	if p.Paid != nil && !*p.Paid {
		p.PaidAt = nil
	}

	if p.Paid != nil && !*p.Paid && p.PaidAt != nil && *p.PaidAt != "" {
		return fmt.Errorf("cannot set paid_at if the payment its not paid")
	}

	if p.Recurrent != nil && *p.Recurrent && p.Frequency != nil && *p.Frequency == "" {
		return fmt.Errorf("frequency is required for recurrent payments")
	}

	if p.Date != nil && *p.Date != "" && p.DueDate != nil && *p.DueDate != "" {
		date, err1 := time.Parse("2006-01-02", *p.Date)
		due_date, err2 := time.Parse("2006-01-02", *p.DueDate)
		if err1 == nil && err2 == nil && due_date.Before(date) {
			return fmt.Errorf("due_date cannot be before payment date")
		}
	}

	return nil
}
