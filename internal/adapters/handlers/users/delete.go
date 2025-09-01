package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *userHandler) DeleteUser(c *gin.Context) {
	id_str := c.Query("id")

	id, err := strconv.Atoi(id_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing id"})
		return
	}

	err = h.UserService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting user"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "user deleted successfully"})
}
