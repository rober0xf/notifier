package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rober0xf/notifier/internal/models"
	"github.com/rober0xf/notifier/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	CreateUserService(email string, password string) error
	GetUserService(email string) (*models.User, error)
	GetAllUsersService() ([]*models.User, error)
	GetUserFromIDService(id uint) (*models.User, error)
	UpdateUserService(*models.User) (*models.User, error)
	DeleteUserService(id uint) error

	// helper functions
	ValidateToken(token_string string) (uint, error)
	ParseUserFromToken(token_string string) (*models.User, error)
	GenerateToken(userID uint, email string) (string, error)
	ExistsUser(ctx context.Context, credentials utils.LoginRequest) (*models.User, error)
}

type AuthService struct {
	db     *gorm.DB
	jwtKey []byte
}

func NewAuthService(db *gorm.DB, jwtKey []byte) *AuthService {
	return &AuthService{
		db:     db,
		jwtKey: jwtKey,
	}
}

var _ AuthServiceInterface = (*AuthService)(nil)

// user services
func (as *AuthService) CreateUserService(email string, password string) error {
	if email == "" || password == "" {
		return utils.ErrInvalidUserData
	}

	user := models.User{
		Email:    email,
		Password: password,
	}

	_, err := as.RegisterUser(&user)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) RegisterUser(user *models.User) (*models.User, error) {
	var existing_user models.User
	err := as.db.Where("email = ?", user.Email).First(&existing_user).Error

	if err == nil {
		return nil, utils.ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	password, err := utils.Hash_password(user.Password)
	if err != nil {
		return nil, utils.ErrPasswordHashing
	}
	user.Password = password

	if err := as.db.Create(user).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "null") || strings.Contains(err.Error(), "invalid") {
			return nil, utils.ErrInvalidUserData
		}
		return nil, err
	}

	return user, nil
}

func (as *AuthService) GetUserService(email string) (*models.User, error) {
	if email == "" {
		return nil, utils.ErrInvalidUserData
	}

	var user models.User

	err := as.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrUserNotFound
	}

	return &user, nil
}

func (as *AuthService) GetAllUsersService() ([]*models.User, error) {
	var users []*models.User

	if err := as.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (as *AuthService) GetUserFromIDService(id uint) (*models.User, error) {
	var user models.User

	if err := as.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (as *AuthService) UpdateUserService(user *models.User) (*models.User, error) {
	var db_user models.User
	if err := as.db.Where("id = ?", user.ID).First(&db_user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	if user.Email == "" || user.Name == "" || user.Password == "" || user.Username == "" {
		return nil, utils.ErrInvalidUserData
	}

	hashed_password, err := utils.Hash_password(user.Password)
	if err != nil {
		return nil, utils.ErrPasswordHashing
	}
	user.Password = hashed_password

	// update the user's fields using the input_user instance
	if err := as.db.Model(&db_user).Updates(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (as *AuthService) DeleteUserService(id uint) error {
	return fmt.Errorf("format string")
}

// helper functions
func (as *AuthService) ValidateToken(token_string string) (uint, error) {
	token, err := jwt.ParseWithClaims(token_string, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return as.jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*utils.JWTClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token claims")
	}

	return claims.UserID, nil
}

func (as *AuthService) ParseUserFromToken(token_string string) (*models.User, error) {
	var user models.User

	userID, err := as.ValidateToken(token_string)
	if err != nil {
		return nil, err
	}

	if err := as.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, fmt.Errorf("internal error: %w", err)
	}

	return &user, nil
}

func (as *AuthService) GenerateToken(userID uint, email string) (string, error) {
	expiration := time.Now().Add(utils.TokenExpirationHours)
	claims := &utils.JWTClaims{
		Email:  email,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			Audience:  []string{"notifier"},
			Issuer:    "notifier-service",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(as.jwtKey)
}

func (as *AuthService) ExistsUser(ctx context.Context, credentials utils.LoginRequest) (*models.User, error) {
	var user models.User

	// check if the user exists
	if err := as.db.WithContext(ctx).Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("internal error: %w", err.Error())
	}

	// compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}
