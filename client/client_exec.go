package client

import (
	"verdmell/utils"
)

// ClientExec struct
type ClientExec struct {
	Checks utils.StringList `json: "checks"`
}

// Run
func (c *ClientExec) Run() error {
	return nil
}