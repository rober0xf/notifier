package shared

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// custom errors
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserData   = errors.New("invalid user data")
	ErrPasswordHashing   = errors.New("error hashing password")
	ErrUserNotFound      = errors.New("user not found")
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
