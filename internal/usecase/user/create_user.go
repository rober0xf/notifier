package user

import (
	"context"
	"errors"
	"fmt"
	"os"

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
		return nil, err
	}

	existingUser, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, repoErr.ErrNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, domainErr.ErrUserAlreadyExists
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, domainErr.ErrPasswordHashing
	}

	verificationToken, err := token.GenerateVerificationToken(12)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username:              username,
		Email:                 email,
		Password:              hashedPassword,
		Active:                false,
		EmailVerificationHash: verificationToken.Hash,
		Timeout:               verificationToken.Timeout,
	}

	// store the user
	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	if os.Getenv("ENV") != "test" {
		htmlBody := mail.VerificationEmailHTML(user.Email, verificationToken.Token, uc.baseURL)

		if err := uc.emailSender.Send([]string{user.Email}, "Verify account", htmlBody); err != nil {
			return nil, fmt.Errorf("error sending email verification: %w", err)

		}
	}

	return user, nil
}
