package httpnews

import (
	"encoding/json"
	"net/http"

	httpapi "github.com/alexandear/news-api/pkg/httpapi"
)

const (
	ErrorCodeInternalError = "INTERNAL_ERROR"
)

func (s *Server) sendOK(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		s.log.WithError(err).Error("failed to encode payload")
	}
}

func (s *Server) sendNotFoundError(w http.ResponseWriter, err error, errorCode, message string) {
	s.log.WithError(err).Info("sending not found error")

	s.sendError(w, http.StatusNotFound, errorCode, message)
}

func (s *Server) sendBadRequestError(w http.ResponseWriter, err error, errorCode, message string) {
	s.log.WithError(err).Info("sending bad request error")

	s.sendError(w, http.StatusBadRequest, errorCode, message)
}

func (s *Server) sendDefaultError(w http.ResponseWriter, err error, message string) {
	s.log.WithError(err).Warn("sending internal error")

	s.sendError(w, http.StatusInternalServerError, ErrorCodeInternalError, message)
}

// sendErr wraps sending of an error in the Error format, and handling the failure to marshal that.
func (s *Server) sendError(w http.ResponseWriter, statusCode int, errorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	apiErr := httpapi.ErrorResponse{
		Error: httpapi.Error{
			Code:    errorCode,
			Message: message,
		},
	}

	if err := json.NewEncoder(w).Encode(apiErr); err != nil {
		s.log.WithError(err).Error("failed to encode error")
	}
}
