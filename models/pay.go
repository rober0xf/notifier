package models

import "time"

type Payment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	NetAmount   float64   `json:"net_amount"`
	GrossAmount float64   `json:"gross_amount"`
	Deductible  float64   `json:"deductible"`
	Name        string    `gorm:"not null" json:"name"`
	Type        string    `gorm:"not null" json:"type"`
	Date        time.Time `gorm:"not null" json:"date"`
	Recurrent   bool      `gorm:"not null" json:"recurrent"`
	Paid        bool      `gorm:"not null" json:"paid"`
}
