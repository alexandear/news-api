package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type post struct {
	ID        string         `db:"id"`
	Title     string         `db:"title"`
	Content   sql.NullString `db:"content"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func (s *Storage) GetAllPosts(ctx context.Context) error {
	query := `SELECT id, title, content, created_at, updated_at FROM posts`
	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to queryx: %w", err)
	}

	var posts []*post
	for rows.Next() {
		p := &post{}
		if err := rows.StructScan(p); err != nil {
			return fmt.Errorf("failed to scan: %w", err)
		}

		posts = append(posts, p)
	}

	return nil
}
