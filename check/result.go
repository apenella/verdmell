/*
Package check is used by verdmell to manage the monitoring checks defined by user
*/
package check

import (
	"time"
	"verdmell/utils"
)

// MetadataResult has meta information about a result
type MetadataResult struct {
	Timestamp   int64         `json:"timestamp"`
	InitTime    time.Time     `json:"inittime"`
	ElapsedTime time.Duration `json:"elapsedtime"`
}

// String transform a MetadataResult object to a string
func (m *MetadataResult) String() string {
	var str string
	var err error

	str, err = utils.ObjectToJSONString(m)
	if err != nil {
		return err.Error()
	}

	return str
}

// Result defines the command execution result
type Result struct {
	Metadata *MetadataResult `json:"metadata"`
	Check    string          `json:"check"`
	Command  string          `json:"command"`
	Output   string          `json:"output"`
	ExitCode int             `json:"exit"`
}

// String transform a Result object to a string
func (r *Result) String() string {
	var str string
	var err error

	str, err = utils.ObjectToJSONString(r)
	if err != nil {
		return err.Error()
	}

	return str
}
