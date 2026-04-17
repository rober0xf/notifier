package entity

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	GoogleID     string    `json:"google_id,omitempty"`
	Name         string    `json:"name,omitempty"`
	IsActive     bool      `json:"is_active"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
