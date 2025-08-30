package storage

import (
	"errors"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/rober0xf/notifier/internal/ports"
	"gorm.io/gorm"
)

var _ ports.AuthRepository = (*Repository)(nil)

func (r *Repository) ExistsUser(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}
