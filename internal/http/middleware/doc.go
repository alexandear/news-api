package middleware

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/go-openapi/runtime/middleware"
)

const (
	docsURL    = "docs"
	swaggerURL = "https://unpkg.com/swagger-ui-dist@4.2.1"
)

// Spec serves swagger spec at basePath/swagger.json.
func Spec(basePath string, rawSpecJSON []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware.Spec(basePath, rawSpecJSON, next)
	}
}

func Doc(serviceName, basePath string) func(http.Handler) http.Handler {
	return DocURL(serviceName, basePath, docsURL)
}

type opts struct {
	SpecURL    string
	SwaggerURL string
	Title      string
}

func DocURL(serviceName, basePath, docsURL string) func(http.Handler) http.Handler {
	if basePath == "" {
		basePath = "/"
	}

	return func(next http.Handler) http.Handler {
		pth := path.Join(basePath, docsURL)
		o := opts{
			SpecURL:    "swagger.json",
			SwaggerURL: swaggerURL,
			Title:      serviceName + " documentation",
		}

		tmpl := template.Must(template.New("swagger-ui").Parse(swaggerTemplate))

		buf := bytes.NewBuffer(nil)
		_ = tmpl.Execute(buf, o)
		b := buf.Bytes()

		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if r.URL.Path == pth {
				rw.Header().Set("Content-Type", "text/html; charset=utf-8")
				rw.WriteHeader(http.StatusOK)

				_, _ = rw.Write(b)
				return
			}

			if next == nil {
				rw.Header().Set("Content-Type", "text/plain")
				rw.WriteHeader(http.StatusNotFound)
				_, _ = rw.Write([]byte(fmt.Sprintf("%q not found", pth)))
				return
			}
			next.ServeHTTP(rw, r)
		})
	}
}

//nolint:lll
const (
	swaggerTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <title>{{ .Title }}</title>
	<link rel="icon" type="image/png" href="https://storage.googleapis.com/nimses-static-cdn/website/home/favicon/favicon-16x16.png" sizes="16x16">
	<link rel="icon" type="image/png" href="https://storage.googleapis.com/nimses-static-cdn/website/home/favicon/favicon-32x32.png" sizes="32x32">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" type="text/css" href="{{ .SwaggerURL }}/swagger-ui.css" >
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }
      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }
      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="{{ .SwaggerURL }}/swagger-ui-bundle.js"></script>
    <script>

    window.onload = function() {
      // Begin Swagger UI call region

      const ui = SwaggerUIBundle({
        url: '{{ .SpecURL }}',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "BaseLayout"
      })
      // End Swagger UI call region
      window.ui = ui
    }
  </script>
  </body>
</html>
`
)
