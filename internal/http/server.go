package httpnews

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	news "github.com/alexandear/news-api/internal"
	httpapi "github.com/alexandear/news-api/pkg/httpapi"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml ../../api/news.openapi.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../../api/news.openapi.yaml

const (
	ErrorCodePostNotFound  = "POST_NOT_FOUND"
	ErrorCodeInvalidPostID = "INVALID_POST_ID"
	ErrorCodeInvalidBody   = "INVALID_BODY"
)

type Storage interface {
	CreatePost(ctx context.Context, params news.PostParams) (news.PostMetadata, error)
	GetPost(ctx context.Context, postID string) (news.Post, error)
	GetAllPosts(ctx context.Context) ([]news.Post, error)
	UpdatePost(ctx context.Context, postID string, params news.PostParams) (news.PostMetadata, error)
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
		s.sendDefaultError(w, err, "failed to get all posts")
		return
	}

	respPosts := make([]httpapi.Post, 0, len(posts))
	for _, p := range posts {
		respPosts = append(respPosts, httpapiPost(p))
	}

	s.sendOK(w, &httpapi.GetAllPostsResponse{
		Posts: respPosts,
	})
}

func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	params, err := s.getPostParams(r.Body)
	if err != nil {
		s.sendBadRequestError(w, err, ErrorCodeInvalidBody, "invalid body")
		return
	}

	meta, err := s.storage.CreatePost(r.Context(), params)
	if err != nil {
		s.sendDefaultError(w, err, "failed to create post")
		return
	}

	s.sendOK(w, httpapiPostMetadata(meta))
}

func (s *Server) getPostParams(body io.Reader) (news.PostParams, error) {
	if body == nil {
		return news.PostParams{}, errors.New("empty body")
	}

	var params news.PostParams
	if err := json.NewDecoder(body).Decode(&params); err != nil {
		return news.PostParams{}, fmt.Errorf("failed to decode body: %w", err)
	}

	if err := params.Validate(); err != nil {
		return news.PostParams{}, err
	}

	return params, nil
}

func (s *Server) DeletePost(w http.ResponseWriter, r *http.Request, postID httpapi.PostID) {
	if err := validatePostID(postID); err != nil {
		s.sendBadRequestError(w, err, ErrorCodeInvalidPostID, "post id is invalid")
		return
	}

	err := s.storage.DeletePost(r.Context(), string(postID))
	switch {
	case errors.Is(err, news.ErrNotFound):
		s.sendNotFoundError(w, err, ErrorCodePostNotFound, "post not found")
	case err != nil:
		s.sendDefaultError(w, err, "failed to delete post")
	default:
		s.sendOK(w, nil)
	}
}

func (s *Server) GetPost(w http.ResponseWriter, r *http.Request, postID httpapi.PostID) {
	if err := validatePostID(postID); err != nil {
		s.sendBadRequestError(w, err, ErrorCodeInvalidPostID, "post id is invalid")
		return
	}

	post, err := s.storage.GetPost(r.Context(), string(postID))
	switch {
	case errors.Is(err, news.ErrNotFound):
		s.sendNotFoundError(w, err, ErrorCodePostNotFound, "post not found")
	case err != nil:
		s.sendDefaultError(w, err, "failed to get post")
	default:
		s.sendOK(w, httpapiPost(post))
	}
}

func (s *Server) UpdatePost(w http.ResponseWriter, r *http.Request, postID httpapi.PostID) {
	if err := validatePostID(postID); err != nil {
		s.sendBadRequestError(w, err, ErrorCodeInvalidPostID, "post id is invalid")
		return
	}

	params, err := s.getPostParams(r.Body)
	if err != nil {
		s.sendBadRequestError(w, err, ErrorCodeInvalidBody, "invalid body")
		return
	}

	pm, err := s.storage.UpdatePost(r.Context(), string(postID), params)
	switch {
	case errors.Is(err, news.ErrNotFound):
		s.sendNotFoundError(w, err, ErrorCodePostNotFound, "post not found")
	case err != nil:
		s.sendDefaultError(w, err, "failed to update post")
	default:
		s.sendOK(w, httpapiPostMetadata(pm))
	}
}

func validatePostID(postID httpapi.PostID) error {
	_, err := uuid.Parse(string(postID))
	return err
}

func httpapiPost(p news.Post) httpapi.Post {
	return httpapi.Post{
		PostData: httpapi.PostData{
			Title:   p.Title,
			Content: &p.Content,
		},
		PostMetadata: httpapiPostMetadata(news.PostMetadata{
			PostID:    p.PostID,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}),
	}
}

func httpapiPostMetadata(p news.PostMetadata) httpapi.PostMetadata {
	return httpapi.PostMetadata{
		PostUpdateMetadata: httpapi.PostUpdateMetadata{
			UpdatedAt: &p.UpdatedAt,
		},
		CreatedAt: &p.CreatedAt,
		Id:        &p.PostID,
	}
}
