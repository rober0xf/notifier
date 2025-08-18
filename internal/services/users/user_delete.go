package users

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"gorm.io/gorm"
)

func (u *Users) Delete(id uint) error {
	var db_user domain.User
	if err := u.db.First(&db_user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrUserNotFound
		}
		return err
	}

	if err := u.db.Delete(&db_user).Error; err != nil {
		return err
	}

	return nil
}
