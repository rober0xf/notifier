package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *userHandler) GetByID(c *gin.Context) {
	id_str := c.Param("id") // comes from the url

	if id_str == "" {
		httphelpers.IDParameterNotProvided(c)
		return
	}

	id, err := strconv.Atoi(id_str)
	if err != nil {
		httphelpers.InvalidIDParameter(c, err)
		return
	}

	user, err := h.UserService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, dto.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching user by id"})
		}
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (h *userHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		httphelpers.EmailParameterNotProvided(c)
		return
	}

	user, err := h.UserService.GetByEmail(email)
	if err != nil {
		if errors.Is(err, dto.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching user by email"})
		}
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (h *userHandler) GetAll(c *gin.Context) {
	users, err := h.UserService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching all users"})
		return
	}
	c.JSON(http.StatusOK, users)
}
