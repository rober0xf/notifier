package users

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"gorm.io/gorm"
)

func (u *Users) Update(user *domain.User) (*domain.User, error) {
	var db_user domain.User
	if err := u.db.Where("id = ?", user.ID).First(&db_user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, err
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		return nil, dto.ErrInvalidUserData
	}

	hashed_password, err := authentication.HashPassword(user.Password)
	if err != nil {
		return nil, dto.ErrPasswordHashing
	}
	user.Password = hashed_password

	// update the user's fields using the input_user instance
	if err := u.db.Model(&db_user).Updates(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
