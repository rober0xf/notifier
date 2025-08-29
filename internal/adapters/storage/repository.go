package storage

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

type AuthRepository struct {
	db     *gorm.DB
	jwtKey []byte
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

func NewAuthRepository(db *gorm.DB, jwtKey []byte) *AuthRepository {
	return &AuthRepository{
		db:     db,
		jwtKey: jwtKey,
	}
}
