package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *userHandler) GetVerificationEmail(c *gin.Context) {
	email := c.Param("email")
	hash_str := c.Param("hash")

	if email == "" || hash_str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification link"})
		return
	}
	user, err := h.UserService.GetVerificationEmail(c, email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully", "user": user})
}
