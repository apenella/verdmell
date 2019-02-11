package client

import (
	"verdmell/check"
	"verdmell/context"
	"verdmell/utils"
)

// ClientExec struct
type Exec struct {
	// Checks
	Checks utils.StringList `json: "checks"`
	// Engine
	Engine *check.CheckEngine `json: "engine"`
	// Context contains information about the runtime state
	Context *context.Context `json: "-"`
}

// Run does the client tasks
func (c *Exec) Run() error {
	// if len(c.Checks) == 0 {
	// 	chk := c.Engine.GetChecks()
	// 	for name,_ := range chk.GetCheck() {
	// 		fmt.Println(name)
	// 	}
	// }

	c.Engine.Start()
	c.Context.Logger.Info(c.Checks)
	//c.Engine.Stop()
	return nil

}
