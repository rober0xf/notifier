package entity

import (
	"time"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	GoogleID     string
	Name         string
	IsActive     bool
	Role         Role
	CreatedAt    time.Time
}
