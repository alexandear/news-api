package http

import (
	"log"
	"net/http"

	httpapi "github.com/alexandear/news-api/pkg/httpapi"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml ../../api/news.openapi.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../../api/news.openapi.yaml

var _ httpapi.ServerInterface = &Server{}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	log.Panicln("implement me")
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
