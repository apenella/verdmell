package main

import (
	"os"
	"bufio"
	"strconv"
	"github.com/apenella/messageOutput"
	"verdmell/environment"
	"verdmell/check"
	"verdmell/service"
)

//
// main
//---------------------------------------------------
func main() {

	var err error
	var env *environment.Environment
	var cks *check.CheckSystem
	var srv *service.ServiceSystem

	exitStatus := 0

	// Call to initialize the Environment
	if err, env = environment.NewEnvironment(); err != nil {
		message.WriteError(err)
		os.Exit(4)
	}
	// get the environment attributes
	//setup := env.GetSetup()
	context := env.GetContext()
	output := env.GetOutput()
	// preparing to destroy the output system
	defer output.DestroyInstance()
	

	// Call to initialize the check system
	if err, cks = check.NewCheckSystem(env); err != nil {
		message.WriteError(err)
		os.Exit(4)
	}
	// Call to initialize the ServiceSystem
	if err,srv = service.NewServiceSystem(env); err != nil {
		message.WriteError(err)
		os.Exit(4)
	}
	// Set the output sample channel for checks as the input's service one
	cks.SetOutputSampleChan(srv.GetInputSampleChan())


	switch(context.ExecutionMode){
	case "cluster":
		env.Output.WriteChInfo("Welcome to Verdmell's server mode. I'm waiting your request on http://"+context.Host+":"+strconv.Itoa(context.Port))
		break
	case "standalone":
		message.Write("")
		message.Write(" # That's Verdmell #\n\tstandalone mode")
		message.Write("")

		checkObj := new(check.CheckObject)
		
		//execute an isolated check
		if context.ExecuteCheck != "" {
			checks := cks.GetChecks()
			if checkObj, err = checks.GetCheckObjectByName(context.ExecuteCheck); err != nil {
				env.Output.WriteChError(err)
				os.Exit(4)	
			}	
			if err := cks.StartCheckSystem(checkObj); err != nil {
				env.Output.WriteChError(err)
				os.Exit(4)
			}
		// execute checks from group
		} else if  context.ExecuteCheckGroup != "" {
			var checks []string
			groups := cks.GetCheckgroups()
			 
			if checks, err = groups.GetCheckgroupByName(context.ExecuteCheckGroup); err != nil {
				env.Output.WriteChError(err)
				os.Exit(4)	
			}
			if err := cks.StartCheckSystem(checks); err != nil {
				env.Output.WriteChError(err)
				os.Exit(4)
			}

		//execute all checks
		} else {
			if err := cks.StartCheckSystem(nil); err != nil {
				env.Output.WriteChError(err)
				os.Exit(4)
			}
		}
	}

	ConsoleReader := bufio.NewReader(os.Stdin)
  message.Write("Press Enter")
 	ConsoleReader.ReadString('\n')

	message.Write(srv.String())
	//message.Write("The status is: "+check.Itoa(exitStatus))
	os.Exit(exitStatus)
}