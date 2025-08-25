package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *userHandler) CreateUser(c *gin.Context) {
	// struct used for decode the input
	var input_user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// check if empty fields
	if input_user.Name == "" || input_user.Email == "" || input_user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty fields"})
		return
	}

	// here we use the service logic
	user, err := h.UserService.Create(input_user.Name, input_user.Email, input_user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
		return
	}

	// we return a custom structure without the password
	c.JSON(http.StatusCreated, gin.H{
		"name":  user.Name,
		"email": user.Email,
	})
}
