package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/usecase/user"
	"github.com/rober0xf/notifier/pkg/auth"
)

func (h *UserHandler) Create(c *gin.Context) {
	var input dto.CreateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, email and password are required. password min length 6"})
		return
	}

	// here we use the service logic
	createdUser, err := h.createUserUC.Execute(c.Request.Context(), input.Username, input.Email, input.Password)
	if err != nil {
		if validationErr, ok := err.(*user.EmailValidationError); ok {
			c.JSON(http.StatusBadRequest,
				gin.H{"error": validationErr.Message,
					"suggestion": validationErr.Suggestion,
				})
			return
		}

		switch {
		case errors.Is(err, domainErr.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		case errors.Is(err, domainErr.ErrPasswordHashing):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		case errors.Is(err, domainErr.ErrValidatingEmail):
			c.JSON(http.StatusBadRequest, gin.H{"error": "error processing email"})
		case errors.Is(err, domainErr.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		case errors.Is(err, domainErr.ErrInvalidUsername):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username data"})
		case errors.Is(err, domainErr.ErrInvalidDomain):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email domain"})
		case errors.Is(err, domainErr.ErrDisposableEmail):
			c.JSON(http.StatusBadRequest, gin.H{"error": "disposable emails are not valid"})
		case errors.Is(err, domainErr.ErrEmailNotReachable):
			c.JSON(http.StatusBadRequest, gin.H{"error": "email is not reachable"})
		case errors.Is(err, domainErr.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be stronger"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	token, err := h.tokenGen.Generate(createdUser.ID, createdUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while token generation"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	auth.SetAuthCookie(c, token, auth.CookieConfig{
		Name:            auth.SessionCookieName,
		ExpirationHours: 160,
		Secure:          false,
		HttpOnly:        true,
	})

	// we return a custom structure without the password
	c.JSON(http.StatusCreated, gin.H{
		"email": createdUser.Email,
		"token": token,
	})
}
