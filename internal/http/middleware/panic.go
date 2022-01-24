package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"

	"github.com/alexandear/news-api/internal/http/debug/stack"
	"github.com/alexandear/news-api/internal/http/httperr"
)

func Panic(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rerr := recover(); rerr != nil {
					defer httperr.SendDefaultError(w, "server panic")

					l := logger.WithError(fmt.Errorf("%v", rerr))

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
