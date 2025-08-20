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
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetAll() ([]domain.User, error)
	GetByID(id uint) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
}
