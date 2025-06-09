package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/rober0xf/notifier/internal/core/shared"
	"golang.org/x/crypto/bcrypt"
)

const (
	TokenExpirationHours = 6
	BearerPrefix         = "BEARER "
	SessionCookieName    = "session_token"
	AuthHeaderName       = "Authorization"
)

func Hash_password(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(pass), err
}

// HTTP handler helper functions
func Is_json_request(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}

func Write_json_response(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Write_error_response(w http.ResponseWriter, status int, message, details string) {
	response := shared.ErrorResponse{Message: message}
	if details != "" {
		log.Printf("Error details: %s", details)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func Set_auth_cookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true, // prevent XSS
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(TokenExpirationHours),
	})
}
