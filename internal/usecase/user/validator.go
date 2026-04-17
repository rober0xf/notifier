package user

import (
	emailverifier "github.com/AfterShip/email-verifier"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const minEntropyBits = 60

type EmailValidationError struct {
	Message    string
	Suggestion string
}

func (e *EmailValidationError) Error() string {
	return e.Message
}

func ValidateEmail(email string, disposableDomains []string) error {
	if disposableDomains == nil {
		disposableDomains = []string{}
	}

	if email == "" {
		return domainErr.ErrInvalidEmailFormat
	}

	verifier := emailverifier.NewVerifier().EnableDomainSuggest()
	verifier = verifier.AddDisposableDomains(disposableDomains)

	// verify the email
	ret, err := verifier.Verify(email)
	if err != nil {
		return domainErr.ErrInvalidEmailFormat
	}

	if !ret.Syntax.Valid {
		return domainErr.ErrInvalidEmailFormat
	}

	// check if can receive emails
	if !ret.HasMxRecords {
		return domainErr.ErrInvalidDomain
	}

	if ret.Disposable {
		return domainErr.ErrDisposableEmail
	}

	if ret.Reachable == "no" {
		return domainErr.ErrEmailNotReachable
	}

	// suggestions for typos
	if ret.Suggestion != "" {
		return &EmailValidationError{
			Message:    "invalid email address",
			Suggestion: ret.Suggestion,
		}
	}

	return nil
}

func ValidatePassword(password string) error {
	err := passwordvalidator.Validate(password, minEntropyBits)
	if err != nil {
		return domainErr.ErrInvalidPassword
	}

	return nil
}
