package storage

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/rober0xf/notifier/internal/ports"
	"gorm.io/gorm"
)

var _ ports.UserRepository = (*Repository)(nil)

func (r *Repository) CreateUser(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.ErrUserAlreadyExists
		}
		return dto.ErrRepository
	}
	return nil
}

func (r *Repository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, dto.ErrRepository
	}
	return &user, nil
}

func (r *Repository) GetAllUsers() ([]domain.User, error) {
	var users []domain.User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, dto.ErrRepository
	}
	if len(users) == 0 {
		return nil, dto.ErrUserNotFound
	}
	return users, nil
}

func (r *Repository) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User

	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_errors.ErrNotFound
		}
		return nil, domain_errors.ErrRepository
	}
	return &user, nil
}

func (r *Repository) UpdateUser(user *domain.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return domain_errors.ErrRepository
	}
	return nil
}

func (r *Repository) DeleteUser(id uint) error {
	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		return domain_errors.ErrRepository
	}
	return nil
}
