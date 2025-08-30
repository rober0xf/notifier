package auth

import (
	"github.com/rober0xf/notifier/internal/ports"
)

type Service struct {
	AuthRepo ports.AuthRepository
	UserRepo ports.UserRepository
	jwtKey   []byte
}

func NewAuthService(repo ports.AuthRepository, jwtKey []byte) *Service {
	return &Service{
		Repo:   repo,
		jwtKey: jwtKey,
	}
}

func (s *Service) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return s.jwtKey, nil
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

func (s *Service) ParseUserFromToken(tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return s.jwtKey, nil
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

func (s *Service) GenerateToken(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"id":    userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtKey)
}

var _ ports.AuthService = (*Service)(nil)
