package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
)

func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id") // comes from the url

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user data"})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive"})
		return
	}

	user, err := h.getUserByIDUC.Execute(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) GetByEmailEmpty(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "email parameter required"})
}

func (h *UserHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := h.getUserByEmailUC.Execute(c.Request.Context(), email)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.getAllUsersUC.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		response = append(response, toUserResponse(&u))
	}

	c.JSON(http.StatusOK, response)
}

func toUserResponse(user *entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Active:   user.Active,
	}
}
