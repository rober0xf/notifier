package entity

import "slices"

type TransactionType string

const (
	TransactionTypeExpense      TransactionType = "expense"
	TransactionTypeIncome       TransactionType = "income"
	TransactionTypeSubscription TransactionType = "subscription"
)

var AllTransactionTypes = []TransactionType{
	TransactionTypeExpense,
	TransactionTypeIncome,
	TransactionTypeSubscription,
}

func (t TransactionType) IsValid() bool {
	return slices.Contains(AllTransactionTypes, t)
}
