package httpnews

import (
	"context"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	news "github.com/alexandear/news-api/internal"
	httpapi "github.com/alexandear/news-api/pkg/httpapi"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml ../../api/news.openapi.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../../api/news.openapi.yaml

const (
	ErrorCodePostNotFound  = "POST_NOT_FOUND"
	ErrorCodeInvalidPostID = "INVALID_POST_ID"
)

type Storage interface {
	CreatePost(ctx context.Context, postID string, params news.CreatePostParams) error
	GetPost(ctx context.Context, postID string) (news.Post, error)
	GetAllPosts(ctx context.Context) ([]news.Post, error)
	UpdatePost(ctx context.Context, postID string, params news.UpdatePostParams) error
	DeletePost(ctx context.Context, postID string) error
}

var _ httpapi.ServerInterface = &Server{}

type Server struct {
	storage Storage

	log *log.Logger
}

func NewServer(log *log.Logger, storage Storage) *Server {
	return &Server{
		log:     log,
		storage: storage,
	}
}

func (s *Server) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := s.storage.GetAllPosts(r.Context())
	if err != nil {
		s.sendDefaultError(w, err, "")
		return
	}

	respPosts := make([]httpapi.Post, 0, len(posts))
	for _, p := range posts {
		respPosts = append(respPosts, httpapiPost(p))
	}

	s.sendOK(w, httpapi.GetAllPostsResponse{
		Posts: respPosts,
	})
}

func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	err := s.storage.CreatePost(r.Context(), "", news.CreatePostParams{})
	switch {
	case errors.Is(err, news.ErrInvalidArgument):
		s.sendBadRequestError(w, err, ErrorCodeInvalidPostID, "invalid create params")
	case err != nil:
		s.sendDefaultError(w, err, "failed to create post")
	default:
		s.sendOK(w, nil)
	}
}

func (s *Server) DeletePost(w http.ResponseWriter, r *http.Request, postID httpapi.PostID) {
	err := s.storage.DeletePost(r.Context(), string(postID))
	switch {
	case errors.Is(err, news.ErrInvalidArgument):
		s.sendBadRequestError(w, err, ErrorCodeInvalidPostID, "post id is invalid")
	case errors.Is(err, news.ErrNotFound):
		s.sendNotFoundError(w, err, ErrorCodePostNotFound, "post not found")
	case err != nil:
		s.sendDefaultError(w, err, "failed to delete post")
	default:
		s.sendOK(w, nil)
	}
}

func (s *Server) GetPost(w http.ResponseWriter, r *http.Request, postID httpapi.PostID) {
	post, err := s.storage.GetPost(r.Context(), string(postID))
	switch {
	case errors.Is(err, news.ErrInvalidArgument):
		s.sendBadRequestError(w, err, ErrorCodeInvalidPostID, "post id is invalid")
	case errors.Is(err, news.ErrNotFound):
		s.sendNotFoundError(w, err, ErrorCodePostNotFound, "post not found")
	case err != nil:
		s.sendDefaultError(w, err, "failed to get post")
	default:
		s.sendOK(w, httpapiPost(post))
	}
}

func (s *Server) UpdatePost(w http.ResponseWriter, r *http.Request, postID httpapi.PostID) {
	err := s.storage.UpdatePost(r.Context(), string(postID), news.UpdatePostParams{})
	switch {
	case errors.Is(err, news.ErrInvalidArgument):
		s.sendBadRequestError(w, err, ErrorCodeInvalidPostID, "invalid update params")
	case errors.Is(err, news.ErrNotFound):
		s.sendNotFoundError(w, err, ErrorCodePostNotFound, "post not found")
	case err != nil:
		s.sendDefaultError(w, err, "failed to update post")
	default:
		s.sendOK(w, nil)
	}
}

func httpapiPost(p news.Post) httpapi.Post {
	return httpapi.Post{
		CreatePostBody: httpapi.CreatePostBody{
			Title:   p.Title,
			Content: &p.Content,
		},
		CreatePostResponse: httpapi.CreatePostResponse{
			UpdatePostResponse: httpapi.UpdatePostResponse{
				UpdatedAt: &p.UpdatedAt,
			},
			CreatedAt: &p.CreatedAt,
			Id:        &p.PostID,
		},
	}
}
