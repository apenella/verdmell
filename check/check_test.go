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
        Timestamp: int64(0),
      },
    },
    {
      desc: "Testing a Check with an invalid interval",
      err: errors.New("(Check::ValidateCheck) Check 'fake_check' has an invalid interval"),
      c: &Check{
        Name: "fake_check",
        Description: "",
        Command: "fake command",
        Depend: []string{},
        ExpirationTime: 0,
        Interval: -1,
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
  NewCheckEngine("")
  t.Log(logger)

  tests := []struct{
      desc string
      c *Check
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
        Timestamp: int64(0),
      },
    },
  }

  for _, test := range tests {
		t.Log(test.desc)
    t.Log(test.c)

		r,_ := test.c.ExecuteCommand()
    t.Log(r)
		// if err != nil && assert.Error(t, err) {
		// 	assert.Equal(t, test.err, err)
		// }
	}
}
