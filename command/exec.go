package command

import (
	"flag"

	"verdmell/agent"
	"verdmell/utils"
)

type ExecCommand struct{}

func (c *ExecCommand) Run(args []string) int {
	// create an agent	
	a := &agent.BasicAgent{} 

	flags := flag.NewFlagSet("exec",flag.ContinueOnError)
	flags.Usage = func() {c.Help()}

	flags.IntVar(&a.Loglevel, "loglevel", 0, "Loglevel definition [0: INFO | 1: WARN | 2: ERROR | 3: DEBUG]")
	flags.Var(&a.Checks,"check","Checks to execute")
	flags.StringVar(&a.Configfile,"configfile","","Configuration file")
	flags.StringVar(&a.Configdir,"configdir","","Folder where configuration is placed")

	if err := flags.Parse(args); err != nil {
		return 1
	}
	
	// start agent
	if err := a.Start(); err != nil {
		return 1
	}
	// get status from agent
	if err := a.Status(); err != nil {
		return 1
	}


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