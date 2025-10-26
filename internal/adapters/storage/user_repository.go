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
	user.Active = false
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

func (r *Repository) GetUserByID(id int) (*domain.User, error) {
	var user domain.User

	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrNotFound
		}
		return nil, dto.ErrRepository
	}
	return &user, nil
}

func (r *Repository) UpdateUser(user *domain.User) error {
	result := r.db.Model(&domain.User{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return dto.ErrRepository
	}
	if result.RowsAffected == 0 {
		return dto.ErrUserNotFound
	}
	return nil
}

func (r *Repository) DeleteUser(id int) error {
	result := r.db.Delete(&domain.User{}, id)
	if result.Error != nil {
		return dto.ErrRepository
	}
	if result.RowsAffected == 0 {
		return dto.ErrUserNotFound
	}
	return nil
}

func (r *Repository) SetActive(user *domain.User) error {
	user.Active = true
	result := r.db.Model(user).Update("active", user.Active)
	if result.Error != nil {
		return dto.ErrRepository
	}
	if result.RowsAffected == 0 {
		return dto.ErrUserNotFound
	}
	return nil
}
