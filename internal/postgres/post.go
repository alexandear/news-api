package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	news "github.com/alexandear/news-api/internal"
)

type post struct {
	ID        string         `db:"id"`
	Title     string         `db:"title"`
	Content   sql.NullString `db:"content"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func (p *post) ToModel() news.Post {
	return news.Post{
		PostID:    p.ID,
		Title:     p.Title,
		Content:   p.Content.String,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (s *Storage) CreatePost(ctx context.Context, postID string, params news.CreatePostParams) error {
	const q = `INSERT INTO posts (id, title, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?) posts`

	if _, err := s.db.ExecContext(ctx, q, postID, params.Title, params.Content, params.CreatedAt,
		params.UpdatedAt); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (s *Storage) GetPost(ctx context.Context, postID string) (news.Post, error) {
	const q = `SELECT id, title, content, created_at, updated_at FROM posts WHERE id = $1`

	var p post
	err := s.db.QueryRowxContext(ctx, q, postID).StructScan(&p)
	if errors.Is(err, sql.ErrNoRows) {
		return news.Post{}, news.ErrNotFound
	}
	if err != nil {
		return news.Post{}, fmt.Errorf("failed to query rowx: %w", err)
	}

	return p.ToModel(), nil
}

func (s *Storage) GetAllPosts(ctx context.Context) ([]news.Post, error) {
	const q = `SELECT id, title, content, created_at, updated_at FROM posts`

	rows, err := s.db.QueryxContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to queryx: %w", err)
	}

	var posts []news.Post
	for rows.Next() {
		p := &post{}
		if err := rows.StructScan(p); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		posts = append(posts, p.ToModel())
	}

	return posts, nil
}

func (s *Storage) UpdatePost(ctx context.Context, postID string, params news.UpdatePostParams) error {
	const qs = `SELECT title, content FROM posts WHERE id = ? FOR UPDATE`
	const qu = `UPDATE posts SET title = ?, content = ?, updated_at = ? WHERE id = ?`

	if err := s.Transaction(ctx, nil, func(ctx context.Context, tx *sqlx.Tx) error {
		var p post
		err := tx.QueryRowx(qs, postID).StructScan(&p)
		if errors.Is(err, sql.ErrNoRows) {
			return news.ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("failed to select: %w", err)
		}

		if _, err := tx.Exec(qu, postID, params.Title, params.Content, params.UpdatedAt); err != nil {
			return fmt.Errorf("failed to update: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

func (s *Storage) DeletePost(ctx context.Context, postID string) error {
	const q = `DELETE FROM posts WHERE id = ?`

	rows, err := s.db.ExecContext(ctx, q, postID)
	if err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if affected == 0 {
		return errors.New("not found")
	}

	return nil
}
