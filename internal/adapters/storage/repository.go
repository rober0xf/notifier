package storage

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
