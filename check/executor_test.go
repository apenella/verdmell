/*
Package check is used by verdmell to manage the monitoring checks defined by user
*/
package check

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test CommandExecutorRun
func TestCommandExecutorRun(t *testing.T) {

	tests := []struct {
		desc string
		c    *Check
		r    *Result
		err  error
	}{
		{
			desc: "Testing command execution",
			err:  nil,
			c: &Check{
				Name:           "test_check",
				Description:    "testing a command",
				Command:        "../test/conf.d/scripts/verdmelltest.sh 0 0 Testing!",
				Depend:         []string{},
				ExpirationTime: 0,
				Interval:       0,
				Timeout:        60,
				Timestamp:      int64(0),
			},
			r: &Result{
				Check:    "",
				Command:  "",
				Output:   "Testing!. (exit: 0)",
				ExitCode: 0,
			},
		},
		{
			desc: "Testing timeout on command execution",
			err:  nil,
			c: &Check{
				Name:           "test_check",
				Description:    "testing a command",
				Command:        "../test/conf.d/scripts/verdmelltest.sh 0 2 Testing!",
				Depend:         []string{},
				ExpirationTime: 0,
				Interval:       0,
				Timeout:        1,
				Timestamp:      int64(0),
			},
			r: &Result{
				Check:    "",
				Command:  "",
				Output:   "The command has not finished after 1 seconds",
				ExitCode: -1,
			},
		},
		{
			desc: "Testing non zero exit code",
			err:  nil,
			c: &Check{
				Name:           "test_check",
				Description:    "testing a command",
				Command:        "../test/conf.d/scripts/verdmelltest.sh 1 0 Testing!",
				Depend:         []string{},
				ExpirationTime: 0,
				Interval:       0,
				Timeout:        60,
				Timestamp:      int64(0),
			},
			r: &Result{
				Check:    "",
				Command:  "",
				Output:   "Testing!. (exit: 1)",
				ExitCode: 1,
			},
		},
		{
			desc: "Testing unexisting command",
			err:  errors.New("(CommandExecutor::Run) exec: \"unexistent\": executable file not found in $PATH"),
			c: &Check{
				Name:           "test_check",
				Description:    "testing echo",
				Command:        "unexistent",
				Depend:         []string{},
				ExpirationTime: 0,
				Interval:       0,
				Timeout:        60,
				Timestamp:      int64(0),
			},
			r: &Result{
				Check:    "",
				Command:  "",
				Output:   "",
				ExitCode: -1,
			},
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		ex := &CommandExecutor{}
		res, err := ex.Run(test.c)

		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err, err.Error())
		} else {
			assert.Equal(t, test.r.Output, res.Output, "Unexpected output")
			assert.Equal(t, test.r.ExitCode, res.ExitCode, "Unexpected exit code")
		}
	}
}
