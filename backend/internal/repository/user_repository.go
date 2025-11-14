package repository

import (
	"github.com/prabalesh/loco/backend/pkg/database"
)

type UserRepositoryPostgres struct {
	db *database.Database
}

func NewUserRepositoryPostgres(db *database.Database) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}
