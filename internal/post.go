package news

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

type Post struct {
	PostID    string
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostMetadata struct {
	PostID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (p *PostParams) Validate() error {
	if len(p.Title) < 3 || len(p.Title) > 50 {
		return fmt.Errorf("title should be more 3 and less 50 chars")
	}

	return nil
}
