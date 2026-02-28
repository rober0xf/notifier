package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/pkg/auth"
)

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
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

	payment, err := h.getPaymentByIDUC.Execute(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	if payment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot access payments from other user"})
		return
	}

	c.JSON(http.StatusOK, toPaymentResponse(payment))
}

func (h *PaymentHandler) GetAllPaymentsFromUser(c *gin.Context) {
	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payments, err := h.getAllPaymentsFromUserUC.Execute(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	response := make([]dto.PaymentResponse, 0, len(payments))
	for _, p := range payments {
		response = append(response, toPaymentResponse(&p))
	}

	c.JSON(http.StatusOK, response)
}

func (h *PaymentHandler) GetAllPayments(c *gin.Context) {
	payments, err := h.getAllPaymentsUC.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := make([]dto.PaymentResponse, 0, len(payments))
	for _, p := range payments {
		response = append(response, toPaymentResponse(&p))
	}

	c.JSON(http.StatusOK, response)
}
