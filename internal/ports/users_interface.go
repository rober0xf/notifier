package ports

import (
	"github.com/rober0xf/notifier/internal/domain"
)

type UserService interface {
	Create(name string, email string, password string) (*domain.User, error)
	Get(email string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	GetUserFromID(id uint) (*domain.User, error)
	Update(*domain.User) (*domain.User, error)
	Delete(id uint) error
}

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id uint) error
}
