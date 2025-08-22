package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/rober0xf/notifier/internal/domain/domain_errors"
	"github.com/rober0xf/notifier/internal/services/mail"
	"golang.org/x/crypto/bcrypt"
)

var TokenExpirationHours = 6

func (s *Service) ValidateToken(token_string string) (uint, error) {
	token, err := jwt.ParseWithClaims(token_string, &dto.JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return s.jwtKey, nil
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
func (s *Service) ParseUserFromToken(token_string string) (*mail.MailSender, error) {
	userID, err := s.ValidateToken(token_string)
	if err != nil {
		return nil, err
	}

	_, err = s.Repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, domain_errors.ErrNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, fmt.Errorf("internal error: %w", err)
	}

	return nil, nil
}

func (s *Service) GenerateToken(userID uint, email string) (string, error) {
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
	return token.SignedString(s.jwtKey)
}

func (s *Service) ExistsUser(ctx context.Context, credentials dto.LoginRequest) (*domain.User, error) {
	var user domain.User

	// check if the user exists
	userPtr, err := s.Repo.GetUserByEmail(credentials.Email)
	if err != nil {
		if errors.Is(err, domain_errors.ErrNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("internal error: %w", err)
	}
	user = *userPtr

	// compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}
