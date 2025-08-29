package payments

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *paymentHandler) GetAllPayments(c *gin.Context) {
	user_id, err := h.Utils.GetUserIDFromRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing user id from request"})
		return
	}

	payments, err := h.PaymentService.GetAllPayments(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching payments"})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *paymentHandler) GetPaymentByID(c *gin.Context) {
	id_str := c.Param("id")
	if id_str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "must provide an id"})
		return
	}

	id, err := strconv.Atoi(id_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id value"})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive"})
		return
	}

	userID, err := h.Utils.GetUserIDFromRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing request"})
		return
	}

	payment, err := h.PaymentService.GetPaymentFromIDAndUserID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching payment"})
		return
	}

	c.JSON(http.StatusOK, payment)
}
