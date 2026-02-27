package http

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/internal/domain/entity"
)

func toPaymentResponse(payment *entity.Payment) dto.PaymentResponse {
	return dto.PaymentResponse{
		ID:         payment.ID,
		Name:       payment.Name,
		Amount:     payment.Amount,
		Type:       payment.Type,
		Category:   payment.Category,
		Date:       payment.Date,
		DueDate:    payment.DueDate,
		Paid:       payment.Paid,
		PaidAt:     payment.PaidAt,
		Recurrent:  payment.Recurrent,
		Frequency:  payment.Frequency,
		ReceiptURL: payment.ReceiptURL,
	}
}

func formatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				return fmt.Sprintf("%s is required", field)
			case "min":
				return fmt.Sprintf("%s must be at least %s chars", field, e.Param())
			case "max":
				return fmt.Sprintf("%s must be at most %s chars", field, e.Param())
			case "gt":
				return fmt.Sprintf("%s must be greater than %s", field, e.Param())
			case "oneof":
				switch field {
				case "type":
					return "type must be: expense, income, or subscription"
				case "category":
					return "category must be: electronics, entertainment, education, clothing, work, or sports"
				case "frequency":
					return "frequency must be: daily, weekly, monthly, or yearly"
				}
			case "datetime":
				return fmt.Sprintf("%s must be in YYYY-MM-DD format", field)
			case "url":
				return fmt.Sprintf("%s must be a valid URL", field)
			}
		}
	}

	return "validation failed"
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

func freqPtrOrNil(f entity.FrequencyType) *entity.FrequencyType {
	if f == "" {
		return nil
	}

	return &f
}
