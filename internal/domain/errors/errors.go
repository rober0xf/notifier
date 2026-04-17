package errors

import "errors"

var (
	// users
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrInvalidUserData       = errors.New("invalid user data")
	ErrUserNotFound          = errors.New("user not found")
	ErrActivating            = errors.New("error activating account")
	ErrInvalidDomain         = errors.New("email domain cannot receive emails")
	ErrEmailNotReachable     = errors.New("email is not reachable")
	ErrInvalidPassword       = errors.New("invalid password, weak strength")
	ErrInvalidEmailFormat    = errors.New("invalid email format")
	ErrDisposableEmail       = errors.New("we do not accept disposables emails")
	ErrAlreadyVerified       = errors.New("user already verified")
	ErrInvalidGoogleID       = errors.New("invalid google_id")
	ErrUserNotVerified       = errors.New("user not verified")
	ErrEmailAlreadyExists    = errors.New("email already in use")
	ErrUsernameAlreadyExists = errors.New("username already in use")
	ErrSendingEmail          = errors.New("error sending email")

	// payments
	ErrPaymentAlreadyExists   = errors.New("payment already exists")
	ErrInvalidPaymentData     = errors.New("invalid payment data")
	ErrPaymentNotFound        = errors.New("payment not found")
	ErrInvalidAmount          = errors.New("amount must be greater than zero")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
	ErrInvalidCategory        = errors.New("invalid category")
	ErrInvalidFrequency       = errors.New("invalid frequency for recurrent payment")
	ErrInvalidDate            = errors.New("date is required")

	// general
	ErrInternalServerError = errors.New("internal server error")
)
