package cmd

import (
	"github.com/namsral/flag"
)

func addPostgresURLFlag(fs *flag.FlagSet, postgresURL *string) {
	fs.StringVar(postgresURL, "postgres_url", "", "Database URL")
}
