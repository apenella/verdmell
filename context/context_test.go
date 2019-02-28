/*
Package context contains execution data details
*/
package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestContextString tests MetadataResult String method
func TestContextString(t *testing.T) {
	tests := []struct {
		desc string
		c    *Context
		s    string
		err  error
	}{
		{
			desc: "Testing Context to string",
			c: &Context{
				Host: "verdmell",
				Port: 5497,
				Cluster: []string{
					"1.1.1.1:5497",
					"2.2.2.2:5497",
				},
				Loglevel:       1,
				ChecksFolder:   "./conf.d/checks",
				ServicesFolder: "./conf.d/services",
			},
			s:   "{\"host\":\"verdmell\",\"port\":5497,\"cluster\":[\"1.1.1.1:5497\",\"2.2.2.2:5497\"],\"loglevel\":1,\"checks\":\"./conf.d/checks\",\"services\":\"./conf.d/services\"}",
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		s := test.c.String()

		assert.Equal(t, test.s, s, "Unexpected output")

	}
}
