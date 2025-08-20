package storage

import (
	"errors"

	"github.com/rober0xf/notifier/internal/domain"
	domainErrors "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/ports"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

var _ ports.UserRepository = (*Repository)(nil)

func (r *Repository) Create(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domainErrors.ErrNotFound
		}
		return domainErrors.ErrRepository
	}
	return nil
}

func (r *Repository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, domainErrors.ErrRepository
	}
	return &user, nil
}

func (r *Repository) GetAll() ([]domain.User, error) {
	var users []domain.User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, domainErrors.ErrRepository
	}
	return users, nil
}

func (r *Repository) GetByID(id uint) (*domain.User, error) {
	var user domain.User

	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, domainErrors.ErrRepository
	}
	return &user, nil
}

func (r *Repository) Update(user *domain.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return domainErrors.ErrRepository
	}
	return nil
}

func (r *Repository) Delete(id uint) error {
	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		return domainErrors.ErrRepository
	}
	return nil
}
