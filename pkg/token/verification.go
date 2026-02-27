package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// email verification token
type VerificationToken struct {
	Token     string // send to user
	Hash      string // store in db
	ExpiresAt time.Time
	Timeout   time.Duration
}

func GenerateVerificationToken(expirationHours int) (*VerificationToken, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	plainToken := hex.EncodeToString(b)
	hash := sha256.Sum256([]byte(plainToken))
	expiresAt := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	return &VerificationToken{
		Token:     plainToken,
		Hash:      hex.EncodeToString(hash[:]),
		ExpiresAt: expiresAt,
		Timeout:   time.Until(expiresAt),
	}, nil
}
