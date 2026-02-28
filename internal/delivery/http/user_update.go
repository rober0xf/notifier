package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/usecase/user"
)

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// parse json to dto
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Username == nil && req.Email == nil && req.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one field is required"})
		return
	}

	// create the user with the new data
	input := user.UpdateUserInput{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	updatedUser, err := h.updateUserUC.Execute(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, domainErr.ErrInvalidUserData):
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive"})
		case errors.Is(err, domainErr.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		case errors.Is(err, domainErr.ErrPasswordHashing):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, toUserResponse(updatedUser))
}
