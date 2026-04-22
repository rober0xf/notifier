package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/pkg/auth"
)

func (h *PaymentHandler) Delete(c *gin.Context) {
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

	role, err := auth.GetRoleFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = h.deletePaymentUC.Execute(c.Request.Context(), id, userID, role)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		case errors.Is(err, auth.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete payment from other user"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to delete payment", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.Status(http.StatusNoContent)
}
