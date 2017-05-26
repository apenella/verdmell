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

// Commands is the mapping of all the available Consul commands.
var Commands map[string]cli.CommandFactory

func init() {
	Commands = map[string]cli.CommandFactory {
		"exec": func() (cli.Command, error) {
			return &command.ExecCommand {}, nil 
		},
		/*
			TODO
			start --> start node as cluster mode
			stop --> stop node
			restart --> restart node

			reload --> reload configuration
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
		Version: "0.0.0",
		HelpFunc: cli.FilteredHelpFunc(included, cli.BasicHelpFunc("verdmell")),
	}


	exitStatus, err = c.Run()
	if err != nil {
		message.WriteError(args)
	}

	os.Exit(exitStatus)
}