package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	maxOpenConnections = 10
	maxIdleConnections = 10
	maxConnLifetime    = 1 * time.Minute
	maxConnIdleTime    = 5 * time.Minute
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(postgresURL string) (*Storage, error) {
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")
	sqlxDB.DB.SetMaxOpenConns(maxOpenConnections)
	sqlxDB.DB.SetMaxIdleConns(maxIdleConnections)
	sqlxDB.DB.SetConnMaxLifetime(maxConnLifetime)
	sqlxDB.DB.SetConnMaxIdleTime(maxConnIdleTime)

	if err = sqlxDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &Storage{
		db: sqlxDB,
	}, nil
}

func (s *Storage) Close() error {
	return s.Close()
}

type TxFunc func(ctx context.Context, tx *sqlx.Tx) error

func (s *Storage) Transaction(ctx context.Context, opts *sql.TxOptions, txFn TxFunc) error {
	tx, err := s.db.BeginTxx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin tx, %w", err)
	}

	if err := txFn(ctx, tx); err != nil {
		if e := tx.Rollback(); e != nil {
			return fmt.Errorf("failed to execute tx, %w", err)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx, %w", err)
	}

	return nil
}
