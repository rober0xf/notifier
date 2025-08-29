package payments

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/domain"
)

type json_payment struct {
	NetAmount   float64 `json:"net_amount" binding:"required"`
	GrossAmount float64 `json:"gross_amount"`
	Deductible  float64 `json:"deductible" binding:"required"`
	Name        string  `gorm:"not null" json:"name" binding:"required"`
	Type        string  `gorm:"not null" json:"type" binding:"required"`
	Date        string  `gorm:"not null" json:"date"`
	Recurrent   bool    `gorm:"not null" json:"recurrent"`
	Paid        bool    `gorm:"not null" json:"paid" binding:"required"`
}

func (h *paymentHandler) CreatePayment(c *gin.Context) {
	var input_payment json_payment

	if err := c.ShouldBindJSON(&input_payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, err := h.Utils.GetUserIDFromRequest(c.Request)
	if err != nil || userID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing request"})
		return
	}

	parsed_date, err := time.Parse("02-01-2006", input_payment.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing date, expected: %d-%m-%y"})
		return
	}

	payment := &domain.Payment{
		UserID:      userID,
		NetAmount:   input_payment.NetAmount,
		GrossAmount: input_payment.GrossAmount,
		Deductible:  input_payment.Deductible,
		Name:        input_payment.Name,
		Type:        input_payment.Type,
		Date:        parsed_date,
		Recurrent:   input_payment.Recurrent,
		Paid:        input_payment.Paid,
	}

	payment, err = h.PaymentService.Create(payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating payment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"name": payment.Name, "type": payment.Type, "amount": payment.NetAmount})
}
