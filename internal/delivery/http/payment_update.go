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

func (h *PaymentHandler) Update(c *gin.Context) {
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

	var req dto.UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.getPaymentByIDUC.Execute(c.Request.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		case errors.Is(err, auth.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot update payment from other user"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to get payment by id", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	err = h.updatePaymentUC.Execute(c.Request.Context(), id, req.ToInput())
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to update payment by id", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.Status(http.StatusNoContent)
}
