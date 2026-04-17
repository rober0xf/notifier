package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/rober0xf/notifier/pkg/auth"
	"github.com/rober0xf/notifier/pkg/email"
	mail "github.com/rober0xf/notifier/pkg/email"
	"github.com/rober0xf/notifier/pkg/token"
)

type CreateUserUseCase struct {
	userRepo          repository.UserRepository
	emailSender       mail.EmailSender
	disposableDomains []string
	baseURL           string
}

func NewCreateUserUseCase(userRepo repository.UserRepository,
	emailSender email.EmailSender,
	disposableDomains []string,
	baseURL string) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:          userRepo,
		emailSender:       emailSender,
		disposableDomains: disposableDomains,
		baseURL:           baseURL,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, username string, email string, password string) (*entity.User, error) {
	if err := ValidateEmail(email, uc.disposableDomains); err != nil {
		return nil, err
	}

	if err := ValidatePassword(password); err != nil {
		return nil, domainErr.ErrInvalidPassword
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("CreateUserUC.Execute failed to hash password: %w", err)
	}

	input := &entity.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	createdUser, err := uc.userRepo.CreateUser(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, repoErr.ErrEmailAlreadyExists):
			return nil, domainErr.ErrEmailAlreadyExists
		case errors.Is(err, repoErr.ErrUsernameAlreadyExists):
			return nil, domainErr.ErrUsernameAlreadyExists
		default:
			return nil, fmt.Errorf("CreateUserUC.Execute failed to create user: %w", err)
		}
	}

	// to not block req
	go func(user *entity.User) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := uc.sendVerificationEmail(ctx, user); err != nil {
			slog.ErrorContext(ctx, "failed to send verification email",
				"user_id", user.ID,
				"error", err,
			)
		}
	}(createdUser)

	return createdUser, nil
}

func (uc *CreateUserUseCase) sendVerificationEmail(ctx context.Context, user *entity.User) error {
	tokenData, err := token.GenerateVerificationToken(12)
	if err != nil {
		return fmt.Errorf("CreateUserUC.sendVerificationEmail failed to generate verification token: %w", err)
	}

	hash := sha256.Sum256([]byte(tokenData.Token))
	tokenHash := hex.EncodeToString(hash[:])

	_, err = uc.userRepo.CreateUserToken(ctx, &entity.UserToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		Purpose:   entity.TokenPurposeEmailVerification,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	})
	if err != nil {
		return fmt.Errorf("CreateUserUC.sendVerificationEmail failed to create user token: %w", err)
	}

	body := mail.VerificationEmailHTML(tokenData.Token, uc.baseURL)

	if err := uc.sendEmailWithRetry(ctx, []string{user.Email}, "verify account", body); err != nil {
		slog.ErrorContext(ctx, "failed to send verification email permanently",
			"user_id", user.ID,
			"error", err,
		)
		return fmt.Errorf("CreateUserUC.sendVerificationEmail failed to send email with retry: %w", err)
	}

	return nil
}

func (uc *CreateUserUseCase) sendEmailWithRetry(ctx context.Context, to []string, subject, body string) error {
	var err error

	for i := range 5 {
		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

		err = uc.emailSender.Send(attemptCtx, to, subject, body)
		cancel()

		if err == nil {
			return nil
		}

		select {
		case <-time.After(time.Duration(1<<i) * time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("CreateUserUC.sendEmailWithRetry failed to send email after retries")
}
