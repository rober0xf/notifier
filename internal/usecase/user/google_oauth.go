package user

import (
	"context"
	"errors"
	"fmt"

	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/rober0xf/notifier/pkg/auth"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
)

type GoogleOAuthUseCase struct {
	userRepo repository.UserRepository
	tokenGen auth.TokenGenerator
}

func NewOAuthUseCase(userRepo repository.UserRepository, tokenGen auth.TokenGenerator) *GoogleOAuthUseCase {
	return &GoogleOAuthUseCase{
		userRepo: userRepo,
		tokenGen: tokenGen,
	}
}

func (uc *GoogleOAuthUseCase) Execute(ctx context.Context, googleID, email, name string) (*LoginOutput, error) {
	if googleID == "" {
		return nil, domainErr.ErrInvalidGoogleID
	}

	if email == "" {
		return nil, domainErr.ErrInvalidEmailFormat
	}

	// search for google_id
	user, err := uc.userRepo.GetUserByGoogleID(ctx, googleID)
	if err == nil {
		return uc.buildLoginOutput(user)
	}

	if !errors.Is(err, repoErr.ErrNotFound) {
		return nil, fmt.Errorf("GoogleOAuthUC.Execute failed to get user by google id: %w", err)
	}

	// search for email
	user, err = uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, repoErr.ErrNotFound) {
		return nil, fmt.Errorf("GoogleOAuthUC.Execute failed to get user by email: %w", err)
	}

	// if exists then link account
	if user != nil {
		err = uc.userRepo.UpdateUserGoogleID(ctx, user.ID, googleID)
		if err != nil {
			switch {
			case errors.Is(err, repoErr.ErrNotFound):
				return nil, domainErr.ErrUserNotFound
			case errors.Is(err, repoErr.ErrGoogleExists):
				return nil, auth.ErrGoogleAccountAlreadyLinked
			default:
				return nil, fmt.Errorf("GoogleOAuthUC.Execute failed to link google account: %w", err)
			}
		}

		user.GoogleID = googleID
		return uc.buildLoginOutput(user)
	}

	// if it doesnt exists then create user
	user, err = uc.userRepo.CreateOAuthUser(ctx, email, name, googleID)
	if err != nil {
		return nil, fmt.Errorf("GoogleOAuthUC.Execute failed to create oauth user: %w", err)
	}

	return uc.buildLoginOutput(user)
}

func (uc *GoogleOAuthUseCase) buildLoginOutput(user *entity.User) (*LoginOutput, error) {
	token, err := uc.tokenGen.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("GoogleOAuthUC.buildLoginOutput failed to generate token: %w", err)
	}

	return &LoginOutput{Token: token, User: user}, nil
}
