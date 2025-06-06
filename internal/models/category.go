package models

type Category struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	UserID    uint   `gorm:"not null" json:"user_id"`
	Name      string `gorm:"not null" json:"name"`
	Priority  uint   `gorm:"not null" json:"priority"`
	Recurrent bool   `json:"recurrent"`
	Notify    bool   `json:"notify"`
}
