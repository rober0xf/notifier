package payments

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

type json_payment struct {
	Name       string                 `json:"name" binding:"required,min=3,max=100"`
	Amount     float64                `json:"amount" binding:"required,gt=0"`
	Type       domain.TransactionType `json:"type" binding:"required,oneof=expense income subscription"`
	Category   domain.CategoryType    `json:"category" binding:"required,oneof=electronics entertainment education clothing work sports"`
	Date       string                 `json:"date" binding:"required,datetime=2006-01-02"`
	DueDate    string                 `json:"due_date" binding:"omitempty,datetime=2006-01-02"`
	Paid       bool                   `json:"paid"`
	PaidAt     string                 `json:"paid_at" binding:"omitempty,datetime=2006-01-02"`
	Recurrent  bool                   `json:"recurrent"`
	Frequency  domain.FrequencyType   `json:"frequency" binding:"omitempty,oneof=daily weekly monthly yearly"`
	ReceiptURL string                 `json:"receipt_url" binding:"omitempty,url"`
}

func (h *paymentHandler) CreatePayment(c *gin.Context) {
	var input_payment json_payment

	// return the custom error
	if err := c.ShouldBindJSON(&input_payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": format_validation_error(err)})
		return
	}
	if err := input_payment.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get user_id from the context
	userIDf, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: user_id not in context"})
		return
	}
	userID, ok := userIDf.(int)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid user_id"})
		return
	}

	// some fields are pointers because they are not always needed
	to_pointer := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	payment := &domain.Payment{
		UserID:     userID,
		Amount:     input_payment.Amount,
		Name:       input_payment.Name,
		Type:       input_payment.Type,
		Category:   input_payment.Category,
		Date:       input_payment.Date,
		DueDate:    to_pointer(input_payment.DueDate),
		Paid:       input_payment.Paid,
		PaidAt:     to_pointer(input_payment.PaidAt),
		Recurrent:  input_payment.Recurrent,
		Frequency:  (*domain.FrequencyType)(to_pointer(string(input_payment.Frequency))),
		ReceiptURL: to_pointer(input_payment.ReceiptURL),
	}

	var err error
	payment, err = h.PaymentService.Create(payment)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrPaymentAlreadyExists):
			c.JSON(http.StatusBadRequest, gin.H{"error": "payment already exits"})
		case errors.Is(err, dto.ErrInternalServerError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"name": payment.Name, "type": payment.Type, "category": payment.Category, "amount": payment.Amount})
}
