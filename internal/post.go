package news

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
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

type CreatePostParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (p *CreatePostParams) Validate() error {
	if len(p.Title) < 3 || len(p.Title) > 50 {
		return fmt.Errorf("title should be more 3 and less 50 chars")
	}

	return nil
}

type UpdatePostParams struct {
	Title   string
	Content string
}
