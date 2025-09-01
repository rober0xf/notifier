package auth

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (s *Service) ParseLoginRequest(c *gin.Context) (dto.LoginRequest, error) {
	var credentials dto.LoginRequest

	// take the fields from the request
	if err := c.ShouldBindJSON(&credentials); err != nil {
		return credentials, fmt.Errorf("invalid json: %w", err)
	}

	if credentials.Email == "" || credentials.Password == "" {
		return credentials, errors.New("email and password are empty")
	}

	return credentials, nil
}
