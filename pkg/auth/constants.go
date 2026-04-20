package auth

import "time"

const (
	DefaultBcryptCost = 10
	TokenExpiration   = 6 * time.Hour
	BearerPrefix      = "Bearer "
	SessionCookieName = "session_token"
	AuthHeaderName    = "Authorization"
	RoleAdmin         = "admin"
	RolUser           = "user"
)
