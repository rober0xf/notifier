package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/rober0xf/notifier/pkg/auth"
)

type VerifyEmailUseCase struct {
	userRepo repository.UserRepository
}

func NewVerifyEmailUseCase(userRepo repository.UserRepository) *VerifyEmailUseCase {
	return &VerifyEmailUseCase{
		userRepo: userRepo,
	}
}

func (uc *VerifyEmailUseCase) Execute(ctx context.Context, email string, token string) (*entity.User, error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrUserNotFound
		}

		return nil, err
	}

	if user.Active {
		return nil, domainErr.ErrAlreadyVerified
	}

	// hash the given token
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashString := hex.EncodeToString(tokenHash[:])

	// compare with the stored hash
	if tokenHashString != user.EmailVerificationHash {
		return nil, auth.ErrInvalidToken
	}

	// ensure we activate
	if err := uc.userRepo.UpdateUserActive(ctx, user.ID, true); err != nil {
		return nil, domainErr.ErrActivating
	}
	user.Active = true

	return user, nil
}
