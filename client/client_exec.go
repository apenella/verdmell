package client

import (
	"verdmell/utils"
)

type ClientExec struct {
	Checks utils.StringList `json: "checks"`
}

func (c *ClientExec) Run() error {
	return nil
}
