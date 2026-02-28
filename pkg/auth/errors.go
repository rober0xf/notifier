package auth

import "errors"

var (
	ErrMissingAuthHeader   = errors.New("missing authorization header")
	ErrInvalidHeaderFormat = errors.New("invalid authorization header format")
	ErrInvalidToken        = errors.New("invalid or expired token")
	ErrInvalidClaims       = errors.New("invalid JWT claims")
	ErrInvalidClaimID      = errors.New("invalid claim id")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrNoToken             = errors.New("no token provided")
	ErrMalformedHeader     = errors.New("malformed authorization header")
	ErrNoCookie            = errors.New("session cookie not found")
)
