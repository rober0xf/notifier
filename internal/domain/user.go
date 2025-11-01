package domain

import "time"

type User struct {
	ID                    int           `json:"id"`
	Username              string        `json:"username"`
	Email                 string        `json:"email"`
	Password              string        `json:"password,omitempty"`
	Active                bool          `json:"active"`
	EmailVerificationHash string        `json:"email_verification_hash"`
	CreatedAt             time.Time     `json:"created_at"`
	Timeout               time.Duration `json:"timeout"`
}
