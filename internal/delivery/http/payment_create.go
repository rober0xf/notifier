package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/pkg/auth"
)

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest
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

	// some fields are pointers because they are not always needed
	payment := &entity.Payment{
		UserID:     userID,
		Amount:     req.Amount,
		Name:       req.Name,
		Type:       req.Type,
		Category:   req.Category,
		Date:       req.Date,
		DueDate:    strPtrOrNil(req.DueDate),
		Paid:       req.Paid,
		PaidAt:     strPtrOrNil(req.PaidAt),
		Recurrent:  req.Recurrent,
		Frequency:  freqPtrOrNil(req.Frequency),
		ReceiptURL: strPtrOrNil(req.ReceiptURL),
	}

	createdPayment, err := h.createPaymentUC.Execute(c.Request.Context(), payment)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "payment already exists"})
		case errors.Is(err, domainErr.ErrInvalidPaymentData):
			c.JSON(http.StatusBadRequest, gin.H{"error": "paid_at is required when paid is true"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusCreated, toPaymentResponse(createdPayment))
}
