package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *Service) GenerateToken(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"id":    userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed_token, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}
	return signed_token, nil
}
