package auth

import "github.com/golang-jwt/jwt/v5"

type GoogleUser struct {
	Sub   string
	Email string
	Name  string
}

type Claims struct {
	Email  string `json:"email"`
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type CookieConfig struct {
	Name          string
	MaxAgeSeconds int
	Secure        bool // true in https
	HttpOnly      bool
	SameSite      int
}
