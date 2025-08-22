package storage

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func NewPaymentRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
