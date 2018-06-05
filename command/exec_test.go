
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
	"exec": func() (cli.Command, error) {
		return new(ExecCommand),nil
	},
}

var tests = []testsRun{
	{
		Args: []string {"exec","-loglevel","1"},
		Value: 1,
	},
	{
		Args: []string {"exec"},
		Value: 1,
	},
	{
		Args: []string {"exec","-configdir","../test/conf.d","-check","first"},
		Value: 0,
	},
	{
		Args: []string {"exec","-configdir","../test/conf.d","-check","first","-check","second"},
		Value: 0,
	},
	{
		Args: []string {"exec","-loglevel","3","-configdir","../test/conf.d"},
		Value: 0,
	},
}


func TestRun(t *testing.T){

	var i int = 0
	for _, test := range tests {
		c := &cli.CLI{
			Args: test.Args,
			Commands: Commands,
		}

		code, err := c.Run()
		if err != nil {
			t.Fatalf("(ExecCommand::TestRun)",err)
		}
		if code != test.Value {
			t.Fatalf("(ExecCommand::TestRun) (%v) Error on arguments: , %v",i,test.Args)
		}
		i++
	}
}
