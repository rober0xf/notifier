package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rober0xf/notifier/internal/core/shared"
	"github.com/rober0xf/notifier/internal/models"
)

// import the interface to avoid the
type AuthUtils struct {
	authService shared.AuthServiceInterface
}

// for routes
func NewAuthUtils(authService shared.AuthServiceInterface) *AuthUtils {
	return &AuthUtils{
		authService: authService,
	}
}

// used in auth_handler.go
func (au *AuthUtils) Parse_login_request(r *http.Request) (shared.LoginRequest, error) {
	var credentials shared.LoginRequest
	defer r.Body.Close() // idk if this is necessary

	// if it comes from json
	if Is_json_request(r) {
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

func (au *AuthUtils) validate_token(token_string string) (uint, error) {
	return au.authService.ValidateToken(token_string)
}

// used in auth_handler.go
func (au *AuthUtils) Get_userID_from_request(r *http.Request) (uint, error) {
	token_string := r.Header.Get(AuthHeaderName)
	if token_string == "" {
		return 0, shared.ErrMissingAuthHeader
	}

	// remove the bearer prefix
	if len(token_string) > 7 && strings.ToUpper(token_string[:7]) == BearerPrefix {
		token_string = token_string[7:]
	} else {
		return 0, shared.ErrInvalidHeaderFormat
	}

	userID, err := au.validate_token(token_string)
	if err != nil {
		return 0, shared.ErrInvalidToken
	}

	return userID, nil
}

// creates a new jwt token for the user
func (au *AuthUtils) generate_token(userID uint, email string) (string, error) {
	return au.authService.GenerateToken(userID, email)
}

// validate the user using context instead of using the db directly
func (au *AuthUtils) exists_user(ctx context.Context, credentials shared.LoginRequest) (*models.User, error) {
	return au.authService.ExistsUser(ctx, credentials)
}
