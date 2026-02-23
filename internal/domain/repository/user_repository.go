package repository

import (
	"context"

	"github.com/rober0xf/notifier/internal/domain/entity"
)

// type UserService interface {
// 	Create(ctx context.Context, username string, email string, password string) (*domain.User, error)
// 	GetByEmail(ctx context.Context, email string) (*domain.User, error)
// 	GetAll(ctx context.Context) ([]domain.User, error)
// 	GetByID(ctx context.Context, id int) (*domain.User, error)
// 	GetVerificationEmail(ctx context.Context, email string) (*domain.User, error)
// 	Update(ctx context.Context, user *domain.User) (*domain.User, error)
// 	Delete(ctx context.Context, id int) error
// }

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
