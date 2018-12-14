package command

import (
	"flag"

	"verdmell/agent"
	"verdmell/check"
	"verdmell/client"
	"verdmell/context"
	"verdmell/engine"
	"verdmell/utils"
)

// ExecCommand
type ExecCommand struct{}

// Run starts the agent and set the context and the engines
func (c *ExecCommand) Run(args []string) int {
	flags := flag.NewFlagSet("exec", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	// Data structure to set the engines required by agent
	e := make(map[uint]engine.Engine)

	ctx := &context.Context{}

	// Create check an empty check engine
	ch := &check.CheckEngine{}
	e[engine.CHECK] = ch

	// Create check an empty client engine
	// In that case the client is an ClientExec
	cl := &client.Client{}
	e[engine.CLIENT] = cl

	ce := &client.ClientExec{
		Engine: ch,
	}

	flags.IntVar(&ctx.Loglevel, "loglevel", 0, "Loglevel definition [0: INFO | 1: WARN | 2: ERROR | 3: DEBUG]")
	flags.StringVar(&ctx.Configfile, "configfile", "", "Configuration file")
	flags.StringVar(&ctx.Configdir, "configdir", "", "Folder where configuration is placed")
	flags.Var(&ce.Checks, "check", "Checks to execute")

	cl.Worker = ce

	// Create an agent
	a := &agent.Agent{
		Ctx:     ctx,
		Engines: e,
		RunOrder: []uint{
			engine.CHECK,
			engine.CLIENT,
		},
	}

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// start agent
	if exit, err := a.Start(); err != nil {
		return exit
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
// String method transform the ExecCommand to string
func (c *ExecCommand) String() string {
	var str string
	var err error

	str, err = utils.ObjectToJSONString(c)
	if err != nil {
		return err.Error()
	}

	return str
}
