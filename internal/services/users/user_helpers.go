package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/rober0xf/notifier/internal/services/mail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var TokenExpirationHours = 6

func (u *Users) ValidateToken(token_string string) (uint, error) {
	token, err := jwt.ParseWithClaims(token_string, &dto.JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return u.jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*dto.JWTClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token claims")
	}

	return claims.UserID, nil
}

// TODO: fix return
func (u *Users) ParseUserFromToken(token_string string) (*mail.MailSender, error) {
	var user domain.User

	userID, err := u.ValidateToken(token_string)
	if err != nil {
		return nil, err
	}

	if err := u.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, fmt.Errorf("internal error: %w", err)
	}

	return nil, nil
}

func (u *Users) GenerateToken(userID uint, email string) (string, error) {
	expiration := time.Now().Add(time.Duration(TokenExpirationHours))
	claims := &dto.JWTClaims{
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
	return token.SignedString(u.jwtKey)
}

func (u *Users) ExistsUser(ctx context.Context, credentials dto.LoginRequest) (*domain.User, error) {
	var user domain.User

	// check if the user exists
	if err := u.db.WithContext(ctx).Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("internal error: %w", err)
	}

	// compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}
