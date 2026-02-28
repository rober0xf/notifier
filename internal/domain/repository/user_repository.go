package repository

import (
	"context"

	"github.com/rober0xf/notifier/internal/domain/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	UpdateUserProfile(ctx context.Context, id int, username, email string) error
	UpdateUserPassword(ctx context.Context, id int, password string) error
	UpdateUserActive(ctx context.Context, id int, active bool) error
	DeleteUser(ctx context.Context, id int) error
}
