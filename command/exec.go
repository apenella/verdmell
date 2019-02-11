package command

import (
	"flag"
	"time"

	"verdmell/agent"
	"verdmell/check"
	"verdmell/client"
	"verdmell/context"
	"verdmell/engine"
	"verdmell/utils"

	"github.com/apenella/messageOutput"
)

// ExecCommand
type ExecCommand struct{}

// Run starts the agent and set the context and the engines
func (c *ExecCommand) Run(args []string) int {
	var loglevel int

	flags := flag.NewFlagSet("exec", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	logger := message.GetMessager()

	// Data structure to set the engines required by agent
	e := make(map[uint]engine.Engine)

	ctx := &context.Context{
		Logger: logger,
	}

	// Create check an empty check engine
	ch := &check.CheckEngine{
		Context: ctx,
	}
	e[engine.CHECK] = ch

	// Create check an empty client engine
	// In that case the client is an ClientExec
	cl := &client.Client{
		Context: ctx,
	}
	e[engine.CLIENT] = cl

	ce := &client.Exec{
		Engine:  ch,
		Context: ctx,
	}

	flags.IntVar(&loglevel, "loglevel", 0, "Loglevel definition [0: INFO | 1: WARN | 2: ERROR | 3: DEBUG]")
	flags.StringVar(&ctx.Configfile, "configfile", "", "Configuration file")
	flags.StringVar(&ctx.Configdir, "configdir", "", "Folder where configuration is placed")
	flags.Var(&ce.Checks, "check", "Checks to execute")

	err := flags.Parse(args)
	if err != nil {
		return 1
	}

	ctx.Logger.SetLogLevel(loglevel)
	cl.Worker = ce

	// Create an agent
	a := &agent.Agent{
		Context: ctx,
		Engines: e,
		RunOrder: []uint{
			engine.CHECK,
			engine.CLIENT,
		},
	}

	// start agent
	exit, err := a.Start()
	if err != nil {
		return exit
	}
	time.Sleep(time.Duration(5) * time.Second)
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
