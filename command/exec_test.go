
package command

import (
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

type testsRun struct {
	Args []string
	Value int
}

var ExecCommands = map[string]cli.CommandFactory{
	"exec": func() (cli.Command, error) {
		return new(ExecCommand),nil
	},
}

func TestRun(t *testing.T){
	tests := []struct{
		desc string
		args []string
		value int
		err error
	}{
		{
			desc: "Basic exec",
			args: []string{"exec"},
			value: 2,
			err: nil,
		},
		{
			desc: "Exec loglevel 1",
			args: []string {"exec","-loglevel","1"},
			value: 2,
			err: nil,
		},
		{
			desc: "One check with configuration directory",
			args: []string {"exec","-configdir","../test/conf.d","-check","first"},
			value: 0,
			err: nil,
		},
		{
			desc: "Two check with configuration directory",
			args: []string {"exec","-configdir","../test/conf.d","-check","first","-check","second"},
			value: 0,
			err: nil,
		},
		{
			desc: "Basic with configuration and loglevel 3",
			args: []string {"exec","-loglevel","3","-configdir","../test/conf.d"},
			value: 0,
			err: nil,
		},
	}

	//var i int = 0
	for _, test := range tests {
		t.Log(test.desc)

		c := &cli.CLI{
			Args: test.args,
			Commands: ExecCommands,
		}

		code,_ := c.Run()
		// if err != nil {
		// 	t.Fatalf("(ExecCommand::TestRun)",err)
		// }
		assert.Equal(t, test.value, code)
		// if code != test.Value {
		// 	t.Fatalf("(ExecCommand::TestRun) (%v) Error on arguments: , %v",i,test.Args)
		// }
		//i++
	}
}
