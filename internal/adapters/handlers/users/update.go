package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (h *userHandler) Update(c *gin.Context) {
	id_str := c.Param("id")

	id, err := strconv.Atoi(id_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing id"})
		return
	}

	var input_user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// parse json
	if err := c.ShouldBindJSON(&input_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// create the user with the new data
	user := &domain.User{
		ID:       id,
		Username: input_user.Username,
		Email:    input_user.Email,
		Password: input_user.Password,
	}

	updated_user, err := h.UserService.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrInvalidUserData):
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive & email format must be correct"})
		case errors.Is(err, dto.ErrInternalServerError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		case errors.Is(err, dto.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, dto.ErrPasswordHashing):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// clean to show in response
	updated_user.Password = ""

	c.JSON(http.StatusOK, updated_user)
}
