package users

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *userHandler) Create(c *gin.Context) {
	// struct used for decode the input
	var input_user struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, email and password are required"})
		return
	}

	// here we use the service logic
	user, err := h.UserService.Create(input_user.Username, input_user.Email, input_user.Password)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrUserAlreadyExists):
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		case errors.Is(err, dto.ErrPasswordHashing):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		case errors.Is(err, dto.ErrInternalServerError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// we return a custom structure without the password
	c.JSON(http.StatusCreated, gin.H{
		"name":  user.Username,
		"email": user.Email,
	})
}
