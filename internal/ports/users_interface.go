package ports

import (
	"github.com/rober0xf/notifier/internal/domain"
)

type UserService interface {
	Create(username string, email string, password string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	GetAll() ([]domain.User, error)
	GetByID(id int) (*domain.User, error)
	GetVerificationEmail(email string) (*domain.User, error)
	Update(*domain.User) (*domain.User, error)
	Delete(id int) error
}

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	GetUserByID(id int) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id int) error
	SetActive(user *domain.User) error
}
