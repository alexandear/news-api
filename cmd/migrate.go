package cmd

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrateCmd struct {
	fs *flag.FlagSet

	DatabaseURL string
}

func NewMigrateCmd() *MigrateCmd {
	fs := flag.NewFlagSet("migrate", flag.ContinueOnError)
	cmd := &MigrateCmd{
		fs: fs,
	}

	cmd.fs.StringVar(&cmd.DatabaseURL, "database_url", "", "Database URL")

	return cmd
}

func (c MigrateCmd) Name() string {
	return c.fs.Name()
}

func (c MigrateCmd) Description() string {
	return "Run SQL migrations"
}

func (c MigrateCmd) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c MigrateCmd) Run() error {
	m, err := migrate.New("file://migrations", c.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate: %w", err)
	}

	if err = m.Up(); errors.Is(err, migrate.ErrNoChange) {
		return nil
	} else if err != nil {
		return fmt.Errorf("migrate up failed: %w", err)
	}
	return nil
}
