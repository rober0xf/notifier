package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/usecase/user"
	"github.com/rober0xf/notifier/pkg/auth"
)

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive"})
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

	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, err := auth.GetRoleFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	updatedUser, err := h.updateUserUC.Execute(c.Request.Context(), input, userID, role)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, domainErr.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		case errors.Is(err, domainErr.ErrUsernameAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "username already in use"})
		case errors.Is(err, auth.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot change other user data"})
		case errors.Is(err, domainErr.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		case errors.Is(err, domainErr.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to update user", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, toUserResponse(updatedUser))
}
