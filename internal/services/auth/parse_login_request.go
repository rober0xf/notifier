package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (s *Service) ParseLoginRequest(r *http.Request) (dto.LoginRequest, error) {
	var credentials dto.LoginRequest
	defer r.Body.Close() // idk if this is necessary

	// if it comes from json
	if httphelpers.IsJSONRequest(r) {
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
