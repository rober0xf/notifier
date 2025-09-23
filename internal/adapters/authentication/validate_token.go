package authentication

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func validateTokenString(tokenString string, jwtKey []byte) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*dto.JWTClaims)
	if !ok || !token.Valid {
		return 0, dto.ErrInvalidClaims
	}

	return claims.UserID, nil
}
