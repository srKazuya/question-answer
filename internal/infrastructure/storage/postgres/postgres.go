// Package postgres provides functionality for interacting with a PostgreSQL database.
package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pressly/goose/v3"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrOpenDB    = errors.New("failed to open database")
	ErrMigration = errors.New("failed to run migrations")
	ErrGormOpen  = errors.New("failed to gorm open")
)

type PostgresStorage struct {
	db *gorm.DB
}

func New(cfg Config) (*PostgresStorage, error) {
	const op = "storage.postgres.NewStrorage"

	sqlDB, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrOpenDB, err)
	}

	if err := goose.Up(sqlDB, cfg.MigrationsPath); err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrMigration, err)
	}

	gormDB, err := gorm.Open(gormpg.New(gormpg.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrGormOpen, err)
	}

	return &PostgresStorage{db: gormDB}, nil
}

