/*
Package check is used by verdmell to manage the monitoring checks defined by user
*/
package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test CommandExecutorRun
func TestCommandExecutorRun(t *testing.T) {

	tests := []struct {
		desc string
		c    *Check
		r    *Result
	}{
		{
			desc: "Testing command execution",
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
				Output:   "exec: \"unexistent\": executable file not found in $PATH",
				ExitCode: -1,
			},
		},
	}

	for _, test := range tests {
		resultCallback := make(chan *Result)

		t.Log(test.desc)
		ex := &CommandExecutor{
			Check:          test.c,
			resultCallback: resultCallback,
		}

		go func() {
			ex.Run()
		}()

		res := <-resultCallback

		assert.Equal(t, test.r.Output, res.Output, "Unexpected output")
		assert.Equal(t, test.r.ExitCode, res.ExitCode, "Unexpected exit code")

		close(resultCallback)
	}
}
