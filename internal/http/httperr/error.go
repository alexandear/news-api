package httperr

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	api "github.com/alexandear/news-api/pkg/httpapi"
)

const (
	ErrorCodeInternalError = "INTERNAL_ERROR"
)

func SendNotFoundError(w http.ResponseWriter, errorCode, message string) {
	sendError(w, http.StatusNotFound, errorCode, message)
}

func SendBadRequestError(w http.ResponseWriter, errorCode, message string) {
	sendError(w, http.StatusBadRequest, errorCode, message)
}

func SendDefaultError(w http.ResponseWriter, message string) {
	sendError(w, http.StatusInternalServerError, ErrorCodeInternalError, message)
}

// sendError wraps sending of an error in the Error format, and handling the failure to marshal that.
func sendError(w http.ResponseWriter, statusCode int, errorCode, message string) {
	apiErr := api.ErrorResponse{
		Error: api.Error{
			Code:    errorCode,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(apiErr); err != nil {
		log.WithError(err).Warn("failed to encode error")
	}
}
