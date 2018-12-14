package client

import (
	"fmt"

	"verdmell/check"
	"verdmell/utils"
)

// ClientExec struct
type ClientExec struct {
	Checks utils.StringList   `json: "checks"`
	Engine *check.CheckEngine `json: "engine"`
}

// Run
func (c *ClientExec) Run() error {
	// if len(c.Checks) == 0 {
	// 	chk := c.Engine.GetChecks()
	// 	for name,_ := range chk.GetCheck() {
	// 		fmt.Println(name)
	// 	}
	// }

	fmt.Println(len(c.Checks), c.Checks)
	return nil
}
