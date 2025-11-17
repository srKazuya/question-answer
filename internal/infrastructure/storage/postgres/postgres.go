// Package postgres provides functionality for interacting with a PostgreSQL database.
package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

)

var (
	ErrOpenDB    = errors.New("failed to open database")
	ErrMigration = errors.New("failed to run migrations")
)

type Storage struct {
	db *sql.DB
}

func New(cfg Config) (*Storage, error) {
	const op = "storage.postgres.NewStrorage"

	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrOpenDB, err)
	}

	if err := goose.Up(db, "internal/infrastructure/storage/postgres/migrations"); err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrMigration, err)
	}

	return &Storage{db: db}, nil
}