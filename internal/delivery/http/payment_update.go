package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/usecase/payment"
	"github.com/rober0xf/notifier/pkg/auth"
)

func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
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

	var req dto.UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	existingPayment, err := h.getPaymentByIDUC.Execute(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domainErr.ErrPaymentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if existingPayment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot update payment from other user"})
		return
	}

	input := payment.UpdatePaymentInput{
		Name:       req.Name,
		Amount:     req.Amount,
		Type:       req.Type,
		Category:   req.Category,
		Date:       req.Date,
		DueDate:    req.DueDate,
		Paid:       req.Paid,
		PaidAt:     req.PaidAt,
		Recurrent:  req.Recurrent,
		Frequency:  req.Frequency,
		ReceiptURL: req.ReceiptURL,
	}

	updatedPayment, err := h.updatePaymentUC.Execute(c.Request.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		case errors.Is(err, domainErr.ErrInvalidAmount):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		case errors.Is(err, domainErr.ErrInvalidTransactionType):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction type"})
		case errors.Is(err, domainErr.ErrInvalidCategory):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
		case errors.Is(err, domainErr.ErrInvalidFrequency):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid frequency"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, toPaymentResponse(updatedPayment))
}
