package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/rober0xf/notifier/internal/domain/entity"
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

func (uc *VerifyEmailUseCase) Execute(ctx context.Context, plainToken string) (*entity.User, error) {
	hash := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(hash[:])

	// verify and mark token as used
	token, err := uc.userRepo.VerifyAndConsumeToken(ctx, tokenHash, entity.TokenPurposeEmailVerification)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, auth.ErrInvalidToken
		}

		return nil, err
	}

	return uc.userRepo.UpdateUserIsActiveReturning(ctx, token.UserID, true)
}
