package domain

type Payment struct {
	ID         int32           `json:"id"`
	UserID     int             `json:"user_id"`
	Name       string          `json:"name"`
	Amount     float64         `json:"amount"`
	Type       TransactionType `json:"type"`
	Category   CategoryType    `json:"category"`
	Date       string          `json:"date"`
	DueDate    *string         `json:"due_date"`
	Paid       bool            `json:"paid"`
	PaidAt     *string         `json:"paid_at"`
	Recurrent  bool            `json:"recurrent"`
	Frequency  *FrequencyType  `json:"frequency"`
	ReceiptURL *string         `json:"receipt_url"`
}
