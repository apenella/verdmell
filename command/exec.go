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

var configuration *configuration.Configuration

// ExecCommand
type ExecCommand struct{}

// Run starts the agent and set the context and the engines
func (c *ExecCommand) Run(args []string) int {
	var loglevel int
	var configfile string
	var configdir string
	var checks utils.StringList

	flags := flag.NewFlagSet("exec", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	flags.IntVar(&loglevel, "loglevel", 0, "Loglevel definition [0: INFO | 1: WARN | 2: ERROR | 3: DEBUG]")
	flags.StringVar(&configfile, "configfile", "", "Configuration file")
	flags.StringVar(&configdir, "configdir", "", "Folder where configuration is placed")
	flags.Var(&checks, "check", "Checks to execute")

	err := flags.Parse(args)
	if err != nil {
		return 1
	}

	// generate a configuration
	configuration, err := configuration.NewConfiguration(configfile, configdir)
	if err != nil {
		msg := "(Agent::init) " + err.Error()
		ctx.Logger.Error(msg)
		return 1
	}

	ctx.Logger.SetLogLevel(loglevel)
	ctx.Configdir = configuration.Checks.Folder
	ctx.Configfile = configuration.Services.Folder

	logger := message.GetMessager()

	// Data structure to set the engines required by agent
	e := make(map[uint]engine.Engine)

	ctx := &context.Context{
		Logger:         logger,
		Host:           configuration.IP,
		Port:           configuration.Port,
		ChecksFolder:   configuration.Checks.Folder,
		ServicesFolder: configuration.Services.Folder,
		Cluster:        configuration.Cluster,
	}

	// Create check an empty check engine
	ch := &check.CheckEngine{
		BasicEngine: engine.BasicEngine{
			Context: ctx,
		},
	}
	e[engine.CHECK] = ch

	// Create check an empty client engine
	// In that case the client is an ClientExec
	cl := &client.Client{
		Context: ctx,
		Worker: &client.Exec{
			Engine:  ch,
			Context: ctx,
			Checks:  checks,
		},
	}
	e[engine.CLIENT] = cl

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
