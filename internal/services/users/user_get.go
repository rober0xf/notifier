package users

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"gorm.io/gorm"
)

func (u *Users) Get(email string) (*domain.User, error) {
	if email == "" {
		return nil, dto.ErrInvalidUserData
	}

	var user domain.User
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *Users) GetAllUsers() ([]*domain.User, error) {
	var users []*domain.User

	if err := u.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (u *Users) GetUserFromID(id uint) (*domain.User, error) {
	var user domain.User

	if err := u.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
