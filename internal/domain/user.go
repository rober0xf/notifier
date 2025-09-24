package domain

type User struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"password,omitempty"`
}
