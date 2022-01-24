package postgres

import (
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
