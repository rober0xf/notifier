package users

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	database "github.com/rober0xf/notifier/internal/ports/db"
	"github.com/rober0xf/notifier/internal/services/mail"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

type EmailValidationError struct {
	Message    string
	Suggestion string
}

func (e *EmailValidationError) Error() string {
	return e.Message
}

func ValidateEmail(email string) error {
	verifier := emailverifier.NewVerifier().EnableDomainSuggest()
	disposable_emails := mail.MustDisposableEmail()
	verifier = verifier.AddDisposableDomains(disposable_emails)

	// verify the email
	ret, err := verifier.Verify(email)
	if err != nil {
		return dto.ErrValidatingEmail
	}
	if !ret.Syntax.Valid {
		return dto.ErrInvalidEmailFormat
	}
	// check if can receive emails
	if !ret.HasMxRecords {
		return dto.ErrInvalidDomain
	}
	if ret.Disposable {
		return dto.ErrDisposableEmail
	}
	if ret.Reachable == "no" {
		return dto.ErrEmailNotReachable
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

func validatePassword(passw string) error {
	const min_entropy_bits = 60
	err := passwordvalidator.Validate(passw, min_entropy_bits)
	if err != nil {
		return dto.ErrInvalidPassword
	}

	return nil
}

func (s *Service) Create(ctx context.Context, username string, email string, password string) (*domain.User, error) {
	// check if the user already exists
	exists_user, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, dto.ErrUserNotFound) {
		return nil, err
	}
	if exists_user != nil {
		return nil, dto.ErrUserAlreadyExists
	}

	hashed, err := authentication.HashPassword(password)
	if err != nil {
		return nil, dto.ErrPasswordHashing
	}

	err = ValidateEmail(email)
	if err != nil {
		return nil, err
	}
	err = validatePassword(password)
	if err != nil {
		return nil, err
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	verification_token := hex.EncodeToString(b) // plain token
	verification_token_hash := sha256.Sum256([]byte(verification_token))
	expires_at := time.Now().Add(12 * time.Hour)
	timeout := time.Until(expires_at)

	user := &domain.User{
		Username:              username,
		Email:                 email,
		Password:              hashed,
		Active:                false,
		EmailVerificationHash: hex.EncodeToString(verification_token_hash[:]),
		Timeout:               timeout,
	}

	// store the user
	if err := s.Repo.CreateUser(ctx, user); err != nil {
		if errors.Is(err, dto.ErrAlreadyExists) {
			return nil, dto.ErrUserAlreadyExists
		}
		if errors.Is(err, dto.ErrRepository) {
			return nil, dto.ErrInternalServerError
		}

		return nil, err
	}

	// send verification email
	url := fmt.Sprintf("http://localhost:3000/v1/users/email_verification/%s/%s", user.Email, verification_token)
	HTMLBody := fmt.Sprintf(`
		<html>
			<body>
				<h1>Email verification</h1>
				<p>click the link to verify:</p>
				<a href="%s">verify account</a>
			</body>
		</html>
	`, url)
	ms := mail.NewMailSender(
		database.GetEnvOrFatal("SMTP_HOST"),
		database.GetEnvOrFatal("SMTP_PORT"),
		database.GetEnvOrFatal("SMTP_USERNAME"),
		database.GetEnvOrFatal("SMTP_PASSWORD"),
	)
	if err := mail.SendMail(ms, []string{user.Email}, "verify account", HTMLBody); err != nil {
		return nil, fmt.Errorf("error sending email verification: %w", err)

	}

	return user, nil
}
