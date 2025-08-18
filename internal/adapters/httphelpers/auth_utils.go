package httphelpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/rober0xf/notifier/internal/ports/auth"
)

const (
	AuthHeaderName = "Authorization"
	BearerPrefix   = "BEARER "
)

type AuthHelper struct {
	authService auth.AuthService
}

// for routes
func NewAuthHelper(authService auth.AuthService) *AuthHelper {
	return &AuthHelper{
		authService: authService,
	}
}

// used in auth_handler.go
func (h *AuthHelper) ParseLoginRequest(r *http.Request) (dto.LoginRequest, error) {
	var credentials dto.LoginRequest
	defer r.Body.Close() // idk if this is necessary

	// if it comes from json
	if IsJSONRequest(r) {
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			return credentials, fmt.Errorf("invalid json: %w", err)
		}
	} else {
		// if it comes from a form
		if err := r.ParseForm(); err != nil {
			return credentials, fmt.Errorf("invalid form: %w", err)
		}
		// set the values from the form
		credentials.Email = r.FormValue("email")
		credentials.Password = r.FormValue("password")
	}

	if credentials.Email == "" || credentials.Password == "" {
		return credentials, errors.New("email and password are empty")
	}

	return credentials, nil
}

func (h *AuthHelper) ValidateToken(token_string string) (uint, error) {
	return h.authService.ValidateToken(token_string)
}

func (h *AuthHelper) GetUserIDFromRequest(r *http.Request) (uint, error) {
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

	userID, err := h.ValidateToken(token_string)
	if err != nil {
		return 0, dto.ErrInvalidToken
	}

	return userID, nil
}

// creates a new jwt token for the user
func (h *AuthHelper) GenerateToken(userID uint, email string) (string, error) {
	return h.authService.GenerateToken(userID, email)
}

// validate the user using context instead of using the db directly
func (h *AuthHelper) ExistsUser(ctx context.Context, credentials dto.LoginRequest) (*domain.User, error) {
	return h.authService.ExistsUser(ctx, credentials)
}
