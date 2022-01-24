package news

import (
	"errors"
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

type CreatePostParams struct {
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdatePostParams struct {
	Title     string
	Content   string
	UpdatedAt time.Time
}
