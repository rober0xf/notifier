package dto

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

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
	ErrNoToken             = errors.New("no token provided")
	ErrMalformedHeader     = errors.New("malformed authorization header")

	// general
	ErrInternalServerError = errors.New("internal server error")

	// repository
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidData   = errors.New("invalid data")
	ErrRepository    = errors.New("repository error")
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=256"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type LoginResponse struct {
	Token string `json:"token"`
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type JWTClaims struct {
	Email  string `json:"email"`
	UserID int    `json:"user_id"`
	jwt.RegisteredClaims
}
