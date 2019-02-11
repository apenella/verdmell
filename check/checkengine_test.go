/*
Package check is used by verdmell to manage the monitoring checks defined by user
*/
package check

import (
	"errors"
	"testing"
	//	"verdmell/utils"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveChecksFromFile(t *testing.T) {

	tests := []struct {
		desc   string
		file   string
		checks map[string][]*Check
		err    error
	}{
		{
			desc: "Test checks file",
			file: "../test/conf.d/checks/checks.json",
			checks: map[string][]*Check{
				"checks": []*Check{
					{
						Name:           "first",
						Description:    "The number one",
						Command:        "./conf.d/scripts/random.sh 4 first",
						Depend:         []string{"second"},
						ExpirationTime: 0,
						Interval:       10,
						Timeout:        0,
						Timestamp:      0,
					},
				},
			},
			err: nil,
		},
		{
			desc:   "Test an empty json",
			file:   "../test/conf.d/checks/empty.json",
			checks: map[string][]*Check{},
			err:    nil,
		},
		{
			desc:   "Test unexisting file",
			file:   "../test/conf.d/checks/unexisting.json",
			checks: map[string][]*Check{},
			err:    errors.New("(checkengine::retrieveChecksFromFile) Checks from '../test/conf.d/checks/unexisting.json' could not be retrieved. (utils::loadJSONFile) Error on loading file '../test/conf.d/checks/unexisting.json' open ../test/conf.d/checks/unexisting.json: no such file or directory"),
		},
		{
			desc:   "Test a file with no checks defined",
			file:   "../test/conf.d/checks/nocheck.json",
			checks: map[string][]*Check{},
			err:    errors.New("(checkengine::retrieveChecksFromFile) Checks from '../test/conf.d/checks/nocheck.json' could not be retrieved. json: cannot unmarshal number into Go value of type []*check.Check"),
		},
		{
			desc:   "Test a malformed file",
			file:   "../test/conf.d/checks/malformed.json",
			checks: map[string][]*Check{},
			err:    errors.New("(checkengine::retrieveChecksFromFile) Checks from '../test/conf.d/checks/malformed.json' could not be retrieved. invalid character '}' looking for beginning of object key string"),
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		checks, err := retrieveChecksFromFile(test.file)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		} else {
			if len(checks) > 0 {
				assert.Equal(t, test.checks["checks"][0].String(), checks[0].String())
			}
		}
	}
}

// func TestInit(t *testing.T) {}
// func TestRun(t *testing.T) {}
// func TestStop(t *testing.T) {}
// func TestStatus(t *testing.T) {}
// func TestGetID(t *testing.T) {}
// func TestGetName(t *testing.T) {}
// func TestGetDependencies(t *testing.T) {}
// func TestGetInputChannel(t *testing.T) {}
// func TestGetStatus(t *testing.T) {}
// func TestSetStatus(t *testing.T) {}
