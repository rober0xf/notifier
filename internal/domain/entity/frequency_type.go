package entity

type FrequencyType string

const (
	FrequencyTypeDaily   FrequencyType = "daily"
	FrequencyTypeWeekly  FrequencyType = "weekly"
	FrequencyTypeMonthly FrequencyType = "monthly"
	FrequencyTypeYearly  FrequencyType = "yearly"
)

func (f FrequencyType) IsValid() bool {
	switch f {
	case FrequencyTypeDaily, FrequencyTypeWeekly, FrequencyTypeMonthly, FrequencyTypeYearly:
		return true
	}
	return false
}
