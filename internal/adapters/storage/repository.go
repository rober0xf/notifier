package storage

import (
	database "github.com/rober0xf/notifier/internal/ports/db"
)

type Repository struct {
	queries *database.Queries
}

func NewAuthRepository(db database.DBTX) *Repository {
	return &Repository{
		queries: database.New(db),
	}
}

func NewUserRepository(db database.DBTX) *Repository {
	return &Repository{
		queries: database.New(db),
	}
}
