package payments

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rober0xf/notifier/internal/domain"
)

/*
for create payment
to return an useful error if the json is invalid
*/
func format_validation_error(err error) string {
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

// this is for create payment. json_payment is local to the module
func (p *json_payment) validate() error {
	if p.Paid && p.PaidAt == "" {
		return fmt.Errorf("paid_at is required when payment is marked as paid")
	}
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

// for update payment
func validate_update_payment(payment *domain.UpdatePayment) error {
	if payment.Paid != nil && *payment.Paid {
		if payment.PaidAt == nil || *payment.PaidAt == "" {
			return fmt.Errorf("paid_at is required when payment is marked as paid")
		}
	}
	if payment.Paid != nil && !*payment.Paid && payment.PaidAt != nil && *payment.PaidAt != "" {
		return fmt.Errorf("cannot set paid_at if the payment its not paid")
	}
	if payment.Recurrent != nil && *payment.Recurrent && payment.Frequency != nil && *payment.Frequency == "" {
		return fmt.Errorf("frequency is required for recurrent payments")
	}
	if payment.Date != nil && *payment.Date != "" && payment.DueDate != nil && *payment.DueDate != "" {
		date, err1 := time.Parse("2006-01-02", *payment.Date)
		due_date, err2 := time.Parse("2006-01-02", *payment.DueDate)
		if err1 == nil && err2 == nil && due_date.Before(date) {
			return fmt.Errorf("due_date cannot be before payment date")
		}
	}
	return nil
}
