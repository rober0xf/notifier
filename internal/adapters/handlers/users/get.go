package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *userHandler) GetByID(c *gin.Context) {
	id_str := c.Param("id") // comes from the url
	id, err := strconv.Atoi(id_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user data"})
		return
	}

	user, err := h.UserService.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrInvalidUserData):
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive"})
		case errors.Is(err, dto.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error fetching user"})
		}
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (h *userHandler) GetByEmailEmpty(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "email parameter required"})
}

func (h *userHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := h.UserService.GetByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrInvalidUserData):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		case errors.Is(err, dto.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, dto.ErrRepository):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (h *userHandler) GetAll(c *gin.Context) {
	users, err := h.UserService.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "no users found"})
		case errors.Is(err, dto.ErrRepository):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching all users"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	c.JSON(http.StatusOK, users)
}
