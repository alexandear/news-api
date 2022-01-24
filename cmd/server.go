package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"

	server "github.com/alexandear/news-api/internal/http"
	httpMiddleware "github.com/alexandear/news-api/internal/http/middleware"
	api "github.com/alexandear/news-api/pkg/httpapi"
)

type ServerCmd struct {
	fs *flag.FlagSet

	Host        string
	Port        int
	PostgresURL string

	log    *log.Logger
	server *http.Server
}

func NewServerCmd() *ServerCmd {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	cmd := &ServerCmd{
		fs: fs,
	}

	cmd.fs.StringVar(&cmd.Host, "host", "localhost", "HTTP host")
	cmd.fs.IntVar(&cmd.Port, "port", 8080, "HTTP port")
	addPostgresURLFlag(cmd.fs, &cmd.PostgresURL)

	return cmd
}

func (c *ServerCmd) Name() string {
	return "server"
}

func (c *ServerCmd) Description() string {
	return "Execute REST server"
}

func (c *ServerCmd) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	c.log = log.New()

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
		httpMiddleware.Panic(c.log),
		httpMiddleware.NewStructuredLogger(c.log),
		httpMiddleware.Spec("", rawSpec),
		httpMiddleware.Doc(swagger.Info.Title, ""),
	)

	serv := server.NewServer()
	api.HandlerFromMux(serv, r)

	c.server = &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", c.Host, c.Port),
	}

	return nil
}

func (c *ServerCmd) Run() error {
	c.log.WithField("address", c.server.Addr).Info("starting server at address")

	if err := c.server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	return nil
}
