package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

const (
	BearerPrefix = "BEARER "
)

func (s *Service) GetUserIDFromRequest(r *http.Request) (int, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return 0, dto.ErrMissingAuthHeader
	}

	// remove the bearer prefix
	if strings.HasPrefix(strings.ToUpper(tokenString), BearerPrefix) {
		tokenString = tokenString[len(BearerPrefix):]
	} else {
		return 0, dto.ErrInvalidHeaderFormat
	}

	userID, err := s.ValidateToken(tokenString, s.jwtKey)
	if err != nil {
		log.Printf("token validation failed: %v", err)
		return 0, dto.ErrInvalidToken
	}

	return userID, nil
}
