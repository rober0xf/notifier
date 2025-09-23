package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *userHandler) Delete(c *gin.Context) {
	id_str := c.Param("id")

	id, err := strconv.Atoi(id_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing id"})
		return
	}

	err = h.UserService.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, dto.ErrRepository):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
