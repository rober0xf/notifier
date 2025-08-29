package storage

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
)

func (r *AuthRepository) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return r.jwtKey, nil
	})
	if err != nil || !token.Valid {
		return 0, dto.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, dto.ErrInvalidClaims
	}

	id_float, ok := claims["id"].(float64)
	if !ok {
		return 0, dto.ErrInvalidClaimID
	}

	return uint(id_float), nil
}

func (r *AuthRepository) ParseUserFromToken(tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return r.jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, dto.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, dto.ErrInvalidClaims
	}

	user := &domain.User{
		ID:    uint(claims["id"].(float64)),
		Email: claims["email"].(string),
	}
	return user, nil
}

func (r *AuthRepository) GenerateToken(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"id":    userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(r.jwtKey)
}

func (r *AuthRepository) ExistsUser(ctx context.Context, credentials dto.LoginRequest) (*domain.User, error) {
	// query user by email first
	var user domain.User
	err := r.db.Where("email = ?", credentials.Email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
