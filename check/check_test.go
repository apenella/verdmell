package check

import (
  "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test ValidateChecks
func TestValidateCheck(t *testing.T) {
	tests := []struct{
      desc string
      c *Check
      err error
  }{
    {
      desc: "Testing a Check with no name",
      err: errors.New("(Check::ValidateCheck) Check requires a Name"),
      c: &Check{
        Name: "",
        Description: "",
        Command: "",
        Depend: []string{},
        ExpirationTime: 0,
        Interval: 0,
        Timeout: 0,
        Timestamp: int64(0),
      },
    },
    {
      desc: "Testing a Check with no command",
      err: errors.New("(Check::ValidateCheck) Check 'fake_check' requires a Command"),
      c: &Check{
        Name: "fake_check",
        Description: "",
        Command: "",
        Depend: []string{},
        ExpirationTime: 0,
        Interval: 0,
        Timeout: 0,
        Timestamp: int64(0),
      },
    },
    {
      desc: "Testing a Check with an invalid expiration time",
      err: errors.New("(Check::ValidateCheck) Check 'fake_check' has an invalid expiration time"),
      c: &Check{
        Name: "fake_check",
        Description: "",
        Command: "fake command",
        Depend: []string{},
        ExpirationTime: -1,
        Interval: 0,
        Timeout: 0,
        Timestamp: int64(0),
      },
    },
    {
      desc: "Testing a Check with an invalid interval",
      err: errors.New("(Check::ValidateCheck) Check 'fake_check' has an invalid interval. Interval is lower than 0"),
      c: &Check{
        Name: "fake_check",
        Description: "",
        Command: "fake command",
        Depend: []string{},
        ExpirationTime: 0,
        Interval: -1,
        Timeout: 0,
        Timestamp: int64(0),
      },
    },
    {
      desc: "Testing a Check with an timeout greater than interval",
      err: errors.New("(Check::ValidateCheck) Check 'fake_check' has an invalid interval.  Timeout should not be greater than interval"),
      c: &Check{
        Name: "fake_check",
        Description: "",
        Command: "fake command",
        Depend: []string{},
        ExpirationTime: 0,
        Interval: 2,
        Timeout: 1,
        Timestamp: int64(0),
      },
    },
  }

  for _, test := range tests {
		t.Log(test.desc)

		err := test.c.ValidateCheck()
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		}
	}
}

// test ExecuteCommand
func TestExecuteCommand(t *testing.T) {

  tests := []struct{
      desc string
      c *Check
      r *Result
      err error
  }{
    {
      desc: "Testing a simple echo command",
      err: nil,
      c: &Check{
        Name: "test_check",
        Description: "testing echo",
        Command: "echo \"hola\"",
        Depend: []string{},
        ExpirationTime: 0,
        Interval: 0,
        Timeout: 60,
        Timestamp: int64(0),
      },
      r: &Result{
        Check: "",
        Command: "",
        Output: "verdmell",
        ExitCode: 0,
      },
    },
    {
      desc: "Testing unexisting command",
      err: errors.New("(Check::ExecuteCommand) Error during 'test_check' command execution."),
      c: &Check{
        Name: "test_check",
        Description: "testing echo",
        Command: "unexistent",
        Depend: []string{},
        ExpirationTime: 0,
        Interval: 0,
        Timeout: 60,
        Timestamp: int64(0),
      },
      r: &Result{
        Check: "",
        Command: "",
        Output: "",
        ExitCode: 0,
      },
    },
  }

  for _, test := range tests {
		t.Log(test.desc)

		res, err := test.c.ExecuteCommand()
    t.Log(res, test.r.Output)

    if err != nil && test.r.Output != "" {
      assert.Equal(t, test.r.Output, res.Output, "Output is not equal.")
    }

    if err != nil && assert.Error(t, err) {
		 	assert.Equal(t, test.err, err)
		}
	}
}
