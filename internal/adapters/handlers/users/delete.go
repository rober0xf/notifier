package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
)

func (h *userHandler) Delete(c *gin.Context) {
	id_str := c.Param("id")
	if id_str == "" {
		httphelpers.IDParameterNotProvided(c)
		return
	}

	id, err := strconv.Atoi(id_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing id"})
		return
	}

	err = h.UserService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "user deleted successfully"})
}
