package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"

	server "github.com/alexandear/news-api/internal/http"
	httpMiddleware "github.com/alexandear/news-api/internal/http/middleware"
	"github.com/alexandear/news-api/internal/postgres"
	api "github.com/alexandear/news-api/pkg/httpapi"
)

type ServerCmd struct {
	fs *flag.FlagSet

	Host        string
	Port        int
	PostgresURL string

	server *http.Server

	closers []io.Closer
}

func NewServerCmd() *ServerCmd {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	cmd := ServerCmd{
		fs: fs,
	}

	cmd.fs.StringVar(&cmd.Host, "host", "localhost", "HTTP host")
	cmd.fs.IntVar(&cmd.Port, "port", 8080, "HTTP port")
	addPostgresURLFlag(cmd.fs, &cmd.PostgresURL)

	return &cmd
}

func (c *ServerCmd) Name() string {
	return "server"
}

func (c *ServerCmd) Description() string {
	return "Execute REST server"
}

func (c *ServerCmd) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *ServerCmd) Run() error {
	logger := log.New()

	defer func() {
		for _, cl := range c.closers {
			err := cl.Close()
			if err != nil {
				logger.WithError(err).Warn("failed to close")
			}
		}
	}()

	swagger, err := api.GetSwagger()
	if err != nil {
		return fmt.Errorf("failed to load swagger spec %w", err)
	}

	rawSpec, err := json.Marshal(&swagger)
	if err != nil {
		return fmt.Errorf("failed to marshal swagger: %w", err)
	}

	r := chi.NewRouter()

	r.Use(
		httpMiddleware.Panic(logger),
		httpMiddleware.NewStructuredLogger(logger),
		httpMiddleware.Spec("", rawSpec),
		httpMiddleware.Doc(swagger.Info.Title, ""),
	)

	stor, err := postgres.NewStorage(c.PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to create new storage: %w", err)
	}
	c.closers = append(c.closers, stor)

	serv := server.NewServer(logger, stor)
	api.HandlerFromMux(serv, r)

	c.server = &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", c.Host, c.Port),
	}
	c.closers = append(c.closers, c.server)

	logger.WithField("address", c.server.Addr).Info("starting server at address")

	if err := c.server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	return nil
}

func addPostgresURLFlag(fs *flag.FlagSet, postgresURL *string) {
	fs.StringVar(postgresURL, "postgres_url", "", "Database URL")
}
