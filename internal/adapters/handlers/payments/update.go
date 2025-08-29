package payments

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/domain"
)

func (h *paymentHandler) UpdatePayment(c *gin.Context) {
	var input_payment json_payment

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := c.ShouldBindJSON(&input_payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request (update payment)"})
		return
	}

	userID, err := h.Utils.GetUserIDFromRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	parsed_date, err := time.Parse("02-01-2006", input_payment.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing date, expected: %d-%m-%y"})
		return
	}

	user_payment := &domain.Payment{
		ID:          uint(id),
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

	updated_payment, err := h.PaymentService.Update(user_payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating payment"})
		return
	}

	c.JSON(http.StatusOK, updated_payment)
}
