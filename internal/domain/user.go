package domain

type User struct {
	ID                    int    `gorm:"primaryKey" json:"id"`
	Username              string `gorm:"not null" json:"username"`
	Email                 string `gorm:"not null;unique" json:"email"`
	Password              string `gorm:"not null" json:"password,omitempty"`
	Active                bool   `gorm:"not null" json:"active"`
	EmailVerificationHash string `json:"email_verification_hash"`
	CreatedAt             string `gorm:"not null" json:"created_at"`
	Timeout               string `json:"timeout"`
}
