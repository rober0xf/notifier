package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/pkg/auth"
)

func (h *PaymentHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payment, err := h.getPaymentByIDUC.Execute(c.Request.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		case errors.Is(err, auth.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot access payments from other users"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to get payment by id", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, toPaymentResponse(*payment))
}

func (h *PaymentHandler) GetMyPayments(c *gin.Context) {
	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payments, err := h.getMyPaymentsUC.Execute(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, domainErr.ErrUserNotVerified):
			c.JSON(http.StatusForbidden, gin.H{"error": "user not verified"})
		case errors.Is(err, auth.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot see payments from other users"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to get all payments from user", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	response := make([]dto.PaymentResponse, 0, len(payments))
	for _, p := range payments {
		response = append(response, toPaymentResponse(p))
	}

	c.JSON(http.StatusOK, response)
}

func (h *PaymentHandler) GetAll(c *gin.Context) {
	payments, err := h.getAllPaymentsUC.Execute(c.Request.Context())
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to get all payments", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := make([]dto.PaymentResponse, 0, len(payments))
	for _, p := range payments {
		response = append(response, toPaymentResponse(p))
	}

	c.JSON(http.StatusOK, response)
}
