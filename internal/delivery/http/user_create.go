package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/usecase/user"
)

func (h *UserHandler) Create(c *gin.Context) {
	var payload dto.RegisterPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, email and password are required. password min length 8"})
		return
	}

	// here we use the service logic
	user, err := h.createUserUC.Execute(c.Request.Context(), payload.Username, payload.Email, payload.Password)
	if err != nil {
		handleCreateUserError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "check your email to verify your account",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"active":   user.IsActive,
		},
	})
}

func handleCreateUserError(c *gin.Context, err error) {
	var validationErr *user.EmailValidationError

	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Message, "suggestion": validationErr.Suggestion})
		return
	}

	switch {
	case errors.Is(err, domainErr.ErrEmailAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
	case errors.Is(err, domainErr.ErrUsernameAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": "username already in use"})
	case errors.Is(err, domainErr.ErrInvalidEmailFormat):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
	case errors.Is(err, domainErr.ErrInvalidDomain):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email domain"})
	case errors.Is(err, domainErr.ErrDisposableEmail):
		c.JSON(http.StatusBadRequest, gin.H{"error": "disposable emails are not valid"})
	case errors.Is(err, domainErr.ErrEmailNotReachable):
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not reachable"})
	case errors.Is(err, domainErr.ErrInvalidPassword):
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be stronger"})
	default:
		slog.ErrorContext(c.Request.Context(), "failed to register user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
