package entity

import "time"

type TokenPurpose string

const (
	TokenPurposeEmailVerification TokenPurpose = "email_verification"
	TokenPurposePasswordReset     TokenPurpose = "password_reset"
)

// email verification token
type VerificationToken struct {
	Token     string // send to user
	Hash      string // store in db
	ExpiresAt time.Time
	Timeout   time.Duration
}

// represents user_tokens in the db
type UserToken struct {
	ID        int
	UserID    int
	TokenHash string
	Purpose   TokenPurpose
	Used      bool
	ExpiresAt time.Time
	CreatedAt time.Time
}
