package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rober0xf/notifier/internal/domain/entity"
)

type GoogleUser struct {
	Sub   string
	Email string
	Name  string
}

type Claims struct {
	Email  string      `json:"email"`
	UserID int         `json:"user_id"`
	Role   entity.Role `json:"role"`
	jwt.RegisteredClaims
}

type CookieConfig struct {
	Name            string
	TokenExpiration time.Duration
	Secure          bool // true in https
	HttpOnly        bool
	SameSite        int
}
