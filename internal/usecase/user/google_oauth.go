package user

import (
	"context"
	"errors"

	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
)

type GoogleOAuthUseCase struct {
	userRepo repository.UserRepository
}

func NewOAuthUseCase(userRepo repository.UserRepository) *GoogleOAuthUseCase {
	return &GoogleOAuthUseCase{
		userRepo: userRepo,
	}
}

func (uc *GoogleOAuthUseCase) Execute(ctx context.Context, googleID, email, name string) (*entity.User, error) {
	if googleID == "" || email == "" {
		return nil, domainErr.ErrInvalidGoogleID
	}

	// search for google_id
	user, err := uc.userRepo.GetUserByGoogleID(ctx, googleID)
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, repoErr.ErrNotFound) {
		return nil, domainErr.ErrUserNotFound
	}

	// search for email
	user, err = uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, repoErr.ErrNotFound) {
		return nil, domainErr.ErrInternalServerError
	}

	// if exists then link account
	if user != nil {
		if err = uc.userRepo.UpdateUserGoogleID(ctx, user.ID, googleID); err != nil {
			return nil, domainErr.ErrInternalServerError
		}
		user.GoogleID = googleID
		return user, nil
	}

	// if it doesnt exists then create user
	user, err = uc.userRepo.CreateOAuthUser(ctx, email, name, googleID)
	if err != nil {
		return nil, domainErr.ErrInternalServerError
	}

	return user, nil
}
