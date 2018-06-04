package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

type testsRun struct {
	Args []string
	Value int
}

var Commands = map[string]cli.CommandFactory{
	"healthcheck": func() (cli.Command, error) {
		return new(HealthCheckCommand),nil
	},
}
