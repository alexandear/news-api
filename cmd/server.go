package cmd

import (
	"flag"
)

type ServerCmd struct {
	fs *flag.FlagSet
}

func NewServerCmd() *ServerCmd {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	return &ServerCmd{
		fs: fs,
	}
}

func (c ServerCmd) Name() string {
	return c.fs.Name()
}

func (c ServerCmd) Description() string {
	return "Execute REST server"
}

func (c ServerCmd) Init(args []string) error {
	return nil
}

func (c ServerCmd) Run() error {
	return nil
}