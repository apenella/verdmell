package command

import (
	"flag"

	"verdmell/agent"
	"verdmell/check"
	"verdmell/client"
	"verdmell/engine"
	"verdmell/utils"
)

/*
	ExecCommand
*/
type ExecCommand struct{}

/*
	Run
*/
func (c *ExecCommand) Run(args []string) int {
	flags := flag.NewFlagSet("exec",flag.ContinueOnError)
	flags.Usage = func() {c.Help()}

	// Data structure to set the engines required by agent
	e := make(map[uint]engine.Engine)
	
	// Create check an empty check engine
	ch := &check.CheckEngine{}
	e[engine.CHECK] = ch

	// Create check an empty client engine
	// In that case the client is an ClientExec
	cl := &client.Client{}
	e[engine.CLIENT] = cl

	ce := &client.ClientExec{}
	// Create an agent
	a := &agent.Agent{
		Engines: e,
	}

	flags.IntVar(&a.Loglevel, "loglevel", 0, "Loglevel definition [0: INFO | 1: WARN | 2: ERROR | 3: DEBUG]")
	flags.StringVar(&a.Configfile,"configfile","","Configuration file")
	flags.StringVar(&a.Configdir,"configdir","","Folder where configuration is placed")
	flags.Var(&ce.Checks,"check","Checks to execute")

	cl.Worker = ce

	if err := flags.Parse(args); err != nil {
		return 1
	}
	
	// start agent
	if err := a.Start(); err != nil {
		return 1
	}

	a.Stop()

	return 0
}

func (c *ExecCommand) Help() string {
	return "Usage: verdmell exec [options]"
}

func (c *ExecCommand) Synopsis() string {
	return "Execute checks on isolated mode"
}



//
// Common methods
//---------------------------------------------------------------------

// String method transform the ExecCommand to string
func (c *ExecCommand) String() string {
	if err, str := utils.ObjectToJsonString(c); err != nil{
		return err.Error()
	} else{
		return str
	} 
}