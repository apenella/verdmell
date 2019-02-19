/*
Package check is used by verdmell to manage the monitoring checks defined by user
*/
package check

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMetadataResultString tests MetadataResult String method
func TestMetadataResultString(t *testing.T) {
	tests := []struct {
		desc string
		m    *MetadataResult
		s    string
		err  error
	}{
		{
			desc: "Testing metadata result",
			m: &MetadataResult{
				Timestamp:   int64(0),
				InitTime:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				ElapsedTime: time.Duration(1) * time.Second,
			},
			s:   "{\"timestamp\":0,\"inittime\":\"2009-11-10T23:00:00Z\",\"elapsedtime\":1000000000}",
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		s := test.m.String()

		assert.Equal(t, test.s, s, "Unexpected output")

	}
}

// TestMetadataResultString tests MetadataResult String method
func TestResultString(t *testing.T) {
	tests := []struct {
		desc string
		r    *Result
		s    string
		err  error
	}{
		{
			desc: "Testing metadata result",
			r: &Result{
				Metadata: &MetadataResult{
					Timestamp:   int64(0),
					InitTime:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
					ElapsedTime: time.Duration(1) * time.Second,
				},
				Check:    "check",
				Command:  "command",
				Output:   "output",
				ExitCode: 0,
			},
			s:   "{\"metadata\":{\"timestamp\":0,\"inittime\":\"2009-11-10T23:00:00Z\",\"elapsedtime\":1000000000},\"check\":\"check\",\"command\":\"command\",\"output\":\"output\",\"exit\":0}",
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		s := test.r.String()

		assert.Equal(t, test.s, s, "Unexpected output")

	}
}
