/*
	Verdmell

	Aleix Penella. 2016
*/
package main

import (
	"os"

	"verdmell/command"
	"github.com/apenella/messageOutput"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available Verdmell commands.
var Commands map[string]cli.CommandFactory

func init() {
	Commands = map[string]cli.CommandFactory {
		"exec": func() (cli.Command, error) {
			return &command.ExecCommand {}, nil
		},
		Commands = map[string]cli.CommandFactory {
			"healthcheck": func() (cli.Command, error) {
				return &command.HealthCheckCommand {}, nil
			},
		/*
			TODO
			//standalone mode options
			exec: run a check
			    - options:
			    	-check: set check name.
			    	-configfile: set configuration file.
			    	-configdir: set configuration directory.
			    	-loglevel: set loglevel.
			    	-silence: no output message.
			healthcheck: run a node health check
				- options:
					-configfile: set configuration file.
			    	-configdir: set configuration directory.
			    	-loglevel: set loglevel.
					-service: ask for an specific service.
					-silence: no output message.

			// cluster mode options
			start: start agent daemon
				- options:
					-ip: set IP address
					-port: set port to listen to.
					-name: node name.
					-cluster: list of nodes to join to.
					-configfile: set configuration file.
			    	-configdir: set configuration directory.
			    	-loglevel: set loglevel.
			stop: stop agent daemon
			reload: reload configuration
			  - reload checks
			  - reload services
			  - options:
					-configfile: set configuration file.
			    	-configdir: set configuration directory.
			    	-loglevel: set loglevel.
			status: agent's current status
			  - check engine
			  - service engine
			  - api engine
			  - cluster engine
			  - ui engine
			  - samples engine
			ps: list check that are currently running
			node health: get a node health status
			  - options:
			    -node: set node. Current node if not set.
			    -service: achieve some services.
			    -silence: no output message.
			service health: get a service healt status
			  - options:
			    -node: set node. Current node if not set.
			    -service: achieve some services.
			    -silence: no output message.
		*/
	}
}

//
// main
//---------------------------------------------------
func main() {
	var err error
	//var env *environment.Environment

	exitStatus := 0

	// Filter out the configtest command from the help display
	var included []string
	for command := range Commands {
		included = append(included, command)
	}

	args := os.Args[1:]

	c := &cli.CLI{
		Args: args,
		Commands: Commands,
		Version: "2.0.0",
		HelpFunc: cli.FilteredHelpFunc(included, cli.BasicHelpFunc("verdmell")),
	}

	exitStatus, err = c.Run()
	if err != nil {
		message.WriteError(args)
	}

	os.Exit(exitStatus)
}
