package auth

import "errors"

var (
	// jwt
	ErrMissingAuthHeader   = errors.New("missing authorization header")
	ErrInvalidHeaderFormat = errors.New("invalid authorization header format")
	ErrInvalidToken        = errors.New("invalid or expired token")
	ErrInvalidClaims       = errors.New("invalid JWT claims")
	ErrInvalidClaimID      = errors.New("invalid claim id")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrNoToken             = errors.New("no token provided")
	ErrGeneratingToken     = errors.New("error generating token")
	ErrMalformedHeader     = errors.New("malformed authorization header")
	ErrRoleMissing         = errors.New("user role missing")
	ErrInvalidRole         = errors.New("invalid role type")

	// cookies
	ErrNoCookie       = errors.New("session cookie not found")
	ErrUserIDNotFound = errors.New("user_id not found in context")
	ErrInvalidUserID  = errors.New("invalid user_id type")

	// google oauth
	ErrInvalidEmailToken = errors.New("invalid email in token")
	ErrEmailNotVerified  = errors.New("email not verified by google")
)
