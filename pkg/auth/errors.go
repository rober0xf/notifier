package auth

import "errors"

var (
	// jwt
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoToken            = errors.New("no token provided")
	ErrMalformedHeader    = errors.New("malformed authorization header")
	ErrRoleMissing        = errors.New("user role missing")
	ErrInvalidRole        = errors.New("invalid role type")
	ErrForbidden          = errors.New("forbidden")

	// cookies
	ErrNoCookie       = errors.New("session cookie not found")
	ErrUserIDNotFound = errors.New("user_id not found in context")
	ErrInvalidUserID  = errors.New("invalid user_id type")

	// google oauth
	ErrInvalidEmailToken          = errors.New("invalid email in token")
	ErrEmailNotVerified           = errors.New("email not verified by google")
	ErrGoogleAccountAlreadyLinked = errors.New("google account already linked")
)
