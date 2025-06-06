package types

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenExpirationHours = 6
	BearerPrefix         = "BEARER "
	SessionCookieName    = "session_token"
	AuthHeaderName       = "Authorization"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type UserInfo struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type JWTClaims struct {
	Email  string `json:"email"`
	UserID uint   `json:"user_id"`
	jwt.RegisteredClaims
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode string `json:"code,omitempty"`
}

// auth helper functions
func Is_json_request(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}

func Write_json_response(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Write_error_response(w http.ResponseWriter, status int, message, details string) {
	response := ErrorResponse{Message: message}
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
