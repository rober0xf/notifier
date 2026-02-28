package user

import (
	"net/mail"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/rober0xf/notifier/internal/domain/errors"
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
	verifier := emailverifier.NewVerifier().EnableDomainSuggest()
	verifier = verifier.AddDisposableDomains(disposableDomains)

	// verify the email
	ret, err := verifier.Verify(email)
	if err != nil {
		return errors.ErrValidatingEmail
	}
	if !ret.Syntax.Valid {
		return errors.ErrInvalidEmailFormat
	}
	// check if can receive emails
	if !ret.HasMxRecords {
		return errors.ErrInvalidDomain
	}
	if ret.Disposable {
		return errors.ErrDisposableEmail
	}
	if ret.Reachable == "no" {
		return errors.ErrEmailNotReachable
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

func ValidateEmailFormat(email string) error {
	if email == "" {
		return errors.ErrInvalidEmailFormat
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.ErrInvalidEmailFormat
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.ErrInvalidEmailFormat
	}

	return nil
}

func ValidatePassword(password string) error {
	err := passwordvalidator.Validate(password, minEntropyBits)
	if err != nil {
		return errors.ErrInvalidPassword
	}

	return nil
}
