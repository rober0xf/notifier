package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/pkg/auth"
)

func (h *PaymentHandler) DeletePayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive"})
		return
	}

	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// verify that it belongs to the user_id
	payment, err := h.getPaymentByIDUC.Execute(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domainErr.ErrPaymentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if payment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete payments from other user"})
		return
	}

	err = h.deletePaymentUC.Execute(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.Status(http.StatusNoContent)
}
