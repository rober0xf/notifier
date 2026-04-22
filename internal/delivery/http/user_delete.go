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

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
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

	err = h.deleteUserUC.Execute(c.Request.Context(), id, userID, role)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrInvalidUserData):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user data"})
		case errors.Is(err, auth.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete other user account"})
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to delete user", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.Status(http.StatusNoContent)
}
