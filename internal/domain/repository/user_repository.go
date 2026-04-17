package repository

import (
	"context"

	"github.com/rober0xf/notifier/internal/domain/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	CreateOAuthUser(ctx context.Context, email, name, googleID string) (*entity.User, error)

	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	GetUserByGoogleID(ctx context.Context, googleID string) (*entity.User, error)

	UpdateUserProfile(ctx context.Context, id int, username, email string) error
	UpdateUserPassword(ctx context.Context, id int, password string) error
	UpdateUserIsActiveReturning(ctx context.Context, id int, isActive bool) (*entity.User, error)
	UpdateUserGoogleID(ctx context.Context, userID int, googleID string) error

	DeleteUser(ctx context.Context, id int) error

	CreateUserToken(ctx context.Context, token *entity.UserToken) (*entity.UserToken, error)
	VerifyAndConsumeToken(ctx context.Context, tokenHash string, purpose entity.TokenPurpose) (*entity.UserToken, error)
	GetTokenByHash(ctx context.Context, tokenHash string, purpose entity.TokenPurpose) (*entity.UserToken, error)
	DeleteOldTokens(ctx context.Context) (int64, error)
}
