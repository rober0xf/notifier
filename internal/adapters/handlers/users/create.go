package users

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/services/users"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=4"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *userHandler) Create(c *gin.Context) {
	var input CreateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, email and password are required. password min length 6"})
		return
	}

	err := users.ValidateEmail(input.Email)
	if err != nil {
		if validation_error, ok := err.(*users.EmailValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": validation_error.Message,
				"suggestion": validation_error.Suggestion})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// here we use the service logic
	user, err := h.UserService.Create(c, input.Username, input.Email, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		case errors.Is(err, dto.ErrPasswordHashing):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		case errors.Is(err, dto.ErrValidatingEmail):
			c.JSON(http.StatusBadRequest, gin.H{"error": "error processing email"})
		case errors.Is(err, dto.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		case errors.Is(err, dto.ErrInvalidUsername):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username data"})
		case errors.Is(err, dto.ErrInvalidDomain):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email domain"})
		case errors.Is(err, dto.ErrDisposableEmail):
			c.JSON(http.StatusBadRequest, gin.H{"error": "disposable emails are not valid"})
		case errors.Is(err, dto.ErrEmailNotReachable):
			c.JSON(http.StatusBadRequest, gin.H{"error": "email is not reachable"})
		case errors.Is(err, dto.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be stronger"})
		case errors.Is(err, dto.ErrInternalServerError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	token, err := h.Utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while token generation"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	authentication.SetAuthCookie(c, token)

	// we return a custom structure without the password
	c.JSON(http.StatusCreated, gin.H{
		"email": user.Email,
		"token": token,
	})
}
