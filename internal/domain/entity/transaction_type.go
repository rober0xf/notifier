package entity

type TransactionType string

const (
	TransactionTypeExpense      TransactionType = "expense"
	TransactionTypeIncome       TransactionType = "income"
	TransactionTypeSubscription TransactionType = "subscription"
)

func (t TransactionType) IsValid() bool {
	switch t {
	case TransactionTypeExpense, TransactionTypeIncome, TransactionTypeSubscription:
		return true
	}
	return false
}
