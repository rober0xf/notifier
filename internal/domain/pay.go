package domain

type Payment struct {
	ID         int             `gorm:"primaryKey" json:"id"`
	UserID     int             `gorm:"not null" json:"user_id"`
	Name       string          `gorm:"not null" json:"name"`
	Amount     float64         `gorm:"not null" json:"amount"`
	Type       TransactionType `gorm:"not null" json:"type"`
	Category   CategoryType    `gorm:"not null" json:"category"`
	Date       string          `gorm:"not null" json:"date"`
	DueDate    *string         `json:"due_date"`
	Paid       bool            `gorm:"not null" json:"paid"`
	PaidAt     *string         `json:"paid_at"`
	Recurrent  bool            `gorm:"not null" json:"recurrent"`
	Frequency  *FrequencyType  `json:"frequency"`
	ReceiptURL *string         `json:"receipt_url"`
}
