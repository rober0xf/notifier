package domain

// sort of enums
type TransactionType string
type CategoryType string
type FrequencyType string

const (
	Expense      TransactionType = "expense"
	Income       TransactionType = "income"
	Subscription TransactionType = "subscription"
)

const (
	Electronics   CategoryType = "electronics"
	Entertainment CategoryType = "entertainment"
	Education     CategoryType = "education"
	Clothing      CategoryType = "clothing"
	Work          CategoryType = "work"
	Sports        CategoryType = "sports"
)

const (
	Daily   FrequencyType = "daily"
	Weekly  FrequencyType = "weekly"
	Monthly FrequencyType = "monthly"
	Yearly  FrequencyType = "yearly"
)

type UpdatePayment struct {
	Name       *string          `json:"name,omitempty"`
	Amount     *float64         `json:"amount,omitempty"`
	Type       *TransactionType `json:"type,omitempty"`
	Category   *CategoryType    `json:"category,omitempty"`
	Date       *string          `json:"date,omitempty"`
	DueDate    *string          `json:"due_date,omitempty"`
	Paid       *bool            `json:"paid,omitempty"`
	PaidAt     *string          `json:"paid_at,omitempty"`
	Recurrent  *bool            `json:"recurrent,omitempty"`
	Frequency  *FrequencyType   `json:"frequency,omitempty"`
	ReceiptURL *string          `json:"receipt_url,omitempty"`
}
