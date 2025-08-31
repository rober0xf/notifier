package auth

import (
	"net/http"
	"strings"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

const (
	AuthHeaderName = "Authorization"
	BearerPrefix   = "BEARER "
)

func (s *Service) GetUserIDFromRequest(r *http.Request) (uint, error) {
	tokenString := r.Header.Get("AuthHeaderName")
	if tokenString == "" {
		return 0, dto.ErrMissingAuthHeader
	}

	// remove the bearer prefix
	if len(tokenString) > 7 && strings.ToUpper(tokenString[:7]) == BearerPrefix {
		tokenString = tokenString[7:]
	} else {
		return 0, dto.ErrInvalidHeaderFormat
	}

	userID, err := s.ValidateToken(tokenString, s.jwtKey)
	if err != nil {
		return 0, dto.ErrInvalidToken
	}

	return userID, nil
}
