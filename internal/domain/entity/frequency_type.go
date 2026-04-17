package entity

import "slices"

type FrequencyType string

const (
	FrequencyTypeDaily   FrequencyType = "daily"
	FrequencyTypeWeekly  FrequencyType = "weekly"
	FrequencyTypeMonthly FrequencyType = "monthly"
	FrequencyTypeYearly  FrequencyType = "yearly"
)

var AllFrequencyTypes = []FrequencyType{
	FrequencyTypeDaily,
	FrequencyTypeWeekly,
	FrequencyTypeMonthly,
	FrequencyTypeYearly,
}

func (f FrequencyType) IsValid() bool {
	return slices.Contains(AllFrequencyTypes, f)
}
