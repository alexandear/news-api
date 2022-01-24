package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"

	"github.com/alexandear/news-api/internal/http/debug/stack"
	httpapi "github.com/alexandear/news-api/pkg/httpapi"
)

func Panic(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rerr := recover(); rerr != nil {
					l := logger.WithError(fmt.Errorf("%v", rerr))

					defer func() {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)

						apiErr := httpapi.ErrorResponse{
							Error: httpapi.Error{
								Code:    "INTERNAL_ERROR",
								Message: "server panic",
							},
						}

						if err := json.NewEncoder(w).Encode(apiErr); err != nil {
							l.WithError(err).Error("failed to encode error")
						}
					}()

					req, err := httputil.DumpRequest(r, false)
					if err != nil {
						l.WithError(err).Errorf("failed to dump the request on panic %v", rerr)
						return
					}

					l.WithField("stack", stack.Stack(3)).WithField("request", string(req)).
						Error("panic middleware triggered")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
