package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// custom errors
var (
	// users
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserData   = errors.New("invalid user data")
	ErrPasswordHashing   = errors.New("error hashing password")
	ErrUserNotFound      = errors.New("user not found")

	// payments
	ErrPaymentAlreadyExists = errors.New("payment already exists")
	ErrInvalidPaymentData   = errors.New("invalid payment data")
	ErrPaymentNotFound      = errors.New("payment not found")

	// auth
	ErrMissingAuthHeader   = errors.New("missing authorization header")
	ErrInvalidHeaderFormat = errors.New("invalid authorization header format")
	ErrInvalidToken        = errors.New("invalid or expired token")
	ErrInvalidClaims       = errors.New("invalid JWT claims")
	ErrInvalidClaimID      = errors.New("invalid claim id")
	ErrInvalidCredentials  = errors.New("invalid credentials")

	// general
	ErrInternalServerError = errors.New("internal server error")
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
