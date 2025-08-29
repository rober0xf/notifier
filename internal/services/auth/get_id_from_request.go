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
	token_string := r.Header.Get("AuthHeaderName")
	if token_string == "" {
		return 0, dto.ErrMissingAuthHeader
	}

	// remove the bearer prefix
	if len(token_string) > 7 && strings.ToUpper(token_string[:7]) == BearerPrefix {
		token_string = token_string[7:]
	} else {
		return 0, dto.ErrInvalidHeaderFormat
	}

	userID, err := s.Repo.ValidateToken(token_string)
	if err != nil {
		return 0, dto.ErrInvalidToken
	}

	return userID, nil
}
