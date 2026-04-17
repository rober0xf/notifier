package errors

import "errors"

var (
	ErrNotFound              = errors.New("resource not found")
	ErrAlreadyExists         = errors.New("resource already exists")
	ErrInvalidData           = errors.New("invalid data")
	ErrGoogleExists          = errors.New("google account already exists")
	ErrEmailAlreadyExists    = errors.New("email already in use")
	ErrUsernameAlreadyExists = errors.New("username already in use")
)
