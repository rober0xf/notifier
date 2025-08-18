package users

import (
	"github.com/rober0xf/notifier/internal/domain"
)

type UserService interface {
	Create(name string, email string, password string) error
	Get(email string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	GetUserFromID(id uint) (*domain.User, error)
	Update(*domain.User) (*domain.User, error)
	Delete(id uint) error
}
