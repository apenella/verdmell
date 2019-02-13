/*
Package check is used by verdmell to manage the monitoring checks defined by user
*/
package check

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Result defines the command execution result
type Result struct {
	Check       string `json:"check"`
	Command     string `json:"command"`
	Output      string `json:"output"`
	ExitCode    int `json:"exit"`
	InitTime    time.Time `json:"inittime"`
	ElapsedTime time.Duration `json:"elapsedtime"`
}

// Executor interface defines and element which could be execute to achieve a Result
// The Run command must receive and id name, a command to execute and the timeout in seconds to stop waiting for the command ends
type Executor interface {
	Run(name string, command string, timeout int) (*Result, error)
}

// ExecutorFactory is a type of function that is a factory for commands.
type ExecutorFactory func() (Executor, error)

// CommandExecutor runs shell commands
type CommandExecutor struct{}

// Run executes the command defined on check an return the result
func (e *CommandExecutor) Run(name string, command string, timeout int) (*Result, error) {
	var elapsedTime time.Duration
	cmdDone := make(chan error)
	defer close(cmdDone)

	exitCode := -1
	output := ""
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmdSplitted := strings.SplitN(command, " ", 2)

	args := []string{}
	if len(cmdSplitted) > 1 {
		args = strings.Split(cmdSplitted[1], " ")
	}

	cmd := exec.Command(cmdSplitted[0], args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	timeInit := time.Now()
	err := cmd.Start()
	if err != nil {
		return nil, errors.New("(CommandExecutor::Run) " + err.Error())
	}

	go func() { cmdDone <- cmd.Wait() }()

	select {
	case err := <-cmdDone:
		elapsedTime = time.Since(timeInit)
		output = strings.TrimSuffix(stdout.String(), "\n")

		// exit status code
		if err == nil {
			exitCode = 0
		} else {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					exitCode = status.ExitStatus()
					if exitCode > 2 || exitCode < 0 {
						exitCode = -1
						output = strings.TrimSuffix(stderr.String(), "\n")
					}
				} else {
					exitCode = -1
					output = strings.TrimSuffix(stderr.String(), "\n")
				}
			}
		}

	case <-time.After(time.Duration(timeout) * time.Second):
		// timed out
		elapsedTime = time.Since(timeInit)
		output = "The command has not finished after " + strconv.Itoa(timeout) + " seconds"
		cmd.Process.Kill()
	}

	//Exit codes
	// OK: 0
	// WARN: 1
	// ERROR: 2
	// UNKNOWN: other (-1)
	return &Result{
		Check:       name,
		Command:     command,
		Output:      output,
		ExitCode:    exitCode,
		InitTime:    timeInit,
		ElapsedTime: elapsedTime,
	}, nil

}
