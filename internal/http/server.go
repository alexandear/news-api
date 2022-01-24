package http

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	httpapi "github.com/alexandear/news-api/pkg/httpapi"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml ../../api/news.openapi.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../../api/news.openapi.yaml

var _ httpapi.ServerInterface = &Server{}

type Storage interface {
	GetAllPosts(ctx context.Context) error
}

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
	err := s.storage.GetAllPosts(r.Context())
	if err != nil {
		s.log.WithError(err).Warn()
	}
}

func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	log.Panicln("implement me")
}

func (s *Server) DeletePost(w http.ResponseWriter, r *http.Request, postID httpapi.PostID) {
	log.Panicln("implement me")
}

func (s *Server) GetPost(w http.ResponseWriter, r *http.Request, postID httpapi.PostID) {
	log.Panicln("implement me")
}

func (s *Server) UpdatePost(w http.ResponseWriter, r *http.Request, postID httpapi.PostID) {
	log.Panicln("implement me")
}
