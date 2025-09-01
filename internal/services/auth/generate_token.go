package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

var TokenExpirationHours = 168

func (s *Service) GenerateToken(userID uint, email string) (string, error) {
	expiration := time.Now().Add(time.Duration(TokenExpirationHours) * time.Hour)
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
