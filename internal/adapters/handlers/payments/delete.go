package payments

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *paymentHandler) DeletePayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	err = h.PaymentService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting payment"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "payment deleted successfully"})
}
