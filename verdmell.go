package main

import (
	"os"
	"strconv"
	"github.com/apenella/messageOutput"
	"verdmell/check"
	"verdmell/environment"
)

//
// main
//---------------------------------------------------
func main() {

	var err error
	var env *environment.Environment
	var cks *check.CheckSystem

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
		env.Output.WriteChError(err)
		os.Exit(4)
	}

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
			_,exitStatus = cks.StartCheckSystem(checkObj)
		// execute checks from group
		} else if  context.ExecuteCheckGroup != "" {
			var checks []string
			groups := cks.GetCheckgroups()
			 
			if checks, err = groups.GetCheckgroupByName(context.ExecuteCheckGroup); err != nil {
				env.Output.WriteChError(err)
				os.Exit(4)	
			}
			_,exitStatus = cks.StartCheckSystem(checks)
		//execute all checks
		} else {
			_,exitStatus = cks.StartCheckSystem(nil)
			output.WriteChDebug("The status is: "+check.Itoa(exitStatus))
		}
	}

	message.Write("The status is: "+check.Itoa(exitStatus))
	os.Exit(exitStatus)
}