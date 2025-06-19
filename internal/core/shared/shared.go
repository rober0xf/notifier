package shared

import (
	"context"

	"github.com/rober0xf/notifier/internal/models"
)

type AuthServiceInterface interface {
	CreateUserService(name string, email string, password string) error
	GetUserService(email string) (*models.User, error)
	GetAllUsersService() ([]*models.User, error)
	GetUserFromIDService(id uint) (*models.User, error)
	UpdateUserService(*models.User) (*models.User, error)
	DeleteUserService(id uint) error

	// helper functions
	ValidateToken(token_string string) (uint, error)
	ParseUserFromToken(token_string string) (*models.User, error)
	GenerateToken(userID uint, email string) (string, error)
	ExistsUser(ctx context.Context, credentials LoginRequest) (*models.User, error)
}

type PaymentServiceInterface interface {
	CreatePaymentService(*models.Payment) error
	GetPaymentService(id uint) (*models.Payment, error)
	GetAllPaymentsService(user_id uint) ([]*models.Payment, error)
	GetPaymentFromIDService(id uint, user_id uint) (*models.Payment, error)
	UpdatePaymentService(*models.Payment) (*models.Payment, error)
	DeletePaymentService(id uint) error
}
