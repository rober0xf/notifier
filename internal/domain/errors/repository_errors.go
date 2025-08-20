package errors

import "errors"

var (
	ErrNotFound = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidData = errors.New("invalid data")
	ErrRepository = errors.New("repository error")
)
