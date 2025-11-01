package payments

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *paymentHandler) GetAllPayments(c *gin.Context) {
	payments, err := h.PaymentService.GetAllPayments(c)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrInternalServerError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *paymentHandler) GetAllPaymentsFromUser(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "must provide an email"})
		return
	}

	payments, err := h.PaymentService.GetAllPaymentsFromUser(c, email)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, dto.ErrInvalidPaymentData):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		case errors.Is(err, dto.ErrInternalServerError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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

	payment, err := h.PaymentService.Get(c, id)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		case errors.Is(err, dto.ErrInternalServerError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, payment)
}
