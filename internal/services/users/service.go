package users

import (
	"github.com/rober0xf/notifier/internal/ports/users"
	"gorm.io/gorm"
)

type Users struct {
	db     *gorm.DB
	jwtKey []byte
}

func NewUsers(db *gorm.DB, jwtKey []byte) *Users {
	return &Users{
		db:     db,
		jwtKey: jwtKey,
	}
}

var _ users.UserService = (*Users)(nil)
