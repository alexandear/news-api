package cmd

import (
	"fmt"
)

type Runner interface {
	Name() string
	Description() string
	Init(args []string) error
	Run() error
}

type RootCmd struct {
	commands []Runner
	help     string
}

func NewRoot() *RootCmd {
	commands := []Runner{
		NewServerCmd(),
		NewMigrateCmd(),
	}

	help := usage()
	for _, cmd := range commands {
		help += fmt.Sprintf("  %s\t\t%s\n", cmd.Name(), cmd.Description())
	}

	return &RootCmd{
		commands: commands,
		help:     help,
	}
}

func (c RootCmd) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("specify a command\n%s", c.help)
	}

	commandName := args[0]
	for _, cmd := range c.commands {
		if cmd.Name() != commandName {
			continue
		}

		if err := cmd.Init(args[1:]); err != nil {
			return fmt.Errorf("failed to init command: %w", err)
		}

		return cmd.Run()
	}

	return fmt.Errorf("unknown command=%s\n%s", commandName, c.help)
}

func usage() string {
	return `Usage:
  news-api [command]

Available Commands:
`
}
