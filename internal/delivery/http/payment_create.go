package http

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/pkg/auth"
)

func (h *PaymentHandler) Create(c *gin.Context) {
	var req dto.CreatePaymentRequest
	validationMeta := gin.H{
		"allowed_types":       entity.AllTransactionTypes,
		"allowed_categories":  entity.AllCategoryTypes,
		"allowed_frequencies": entity.AllFrequencyTypes,
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err), "meta": validationMeta})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "meta": validationMeta})
		return
	}

	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	paidAt := strPtrOrNil(req.PaidAt)
	if req.Paid && paidAt == nil {
		now := time.Now().Format("2006-01-02")
		paidAt = &now
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
		PaidAt:     paidAt,
		Recurrent:  req.Recurrent,
		Frequency:  freqPtrOrNil(req.Frequency),
		ReceiptURL: strPtrOrNil(req.ReceiptURL),
	}

	createdPayment, err := h.createPaymentUC.Execute(c.Request.Context(), payment)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "payment already exists"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to create payment", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusCreated, toPaymentResponse(*createdPayment))
}
