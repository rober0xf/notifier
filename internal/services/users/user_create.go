package users

import (
	"errors"
	"strings"

	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"gorm.io/gorm"
)

func (u *Users) Create(name string, email string, password string) error {
	if name == "" || email == "" || password == "" {
		return dto.ErrInvalidUserData
	}

	user := domain.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	_, err := u.Register(&user)
	if err != nil {
		return err
	}

	return nil
}

func (u *Users) Register(user *domain.User) (*domain.User, error) {
	var existing_user domain.User
	err := u.db.Where("email = ?", user.Email).First(&existing_user).Error

	if err == nil {
		return nil, dto.ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	password, err := authentication.HashPassword(user.Password)
	if err != nil {
		return nil, dto.ErrPasswordHashing
	}
	user.Password = password

	if err := u.db.Create(user).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "null") || strings.Contains(err.Error(), "invalid") {
			return nil, dto.ErrInvalidUserData
		}
		return nil, err
	}

	return user, nil
}
