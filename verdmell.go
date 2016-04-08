/*
	Verdmell

	Aleix Penella. 2016
*/
package main

import (
	"os"
	"strconv"
	"github.com/apenella/messageOutput"
	"verdmell/environment"
	"verdmell/sample"
	"verdmell/check"
	"verdmell/service"
	"verdmell/cluster"
	"verdmell/api"
	"verdmell/ui"
)

//
// main
//---------------------------------------------------
func main() {

	var err error
	var env *environment.Environment
	var sam *sample.SampleEngine
	var cks *check.CheckEngine
	var srv *service.ServiceEngine
	var cltr *cluster.ClusterEngine

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

  // Call to initialize cluster engine
	if err, cltr = cluster.NewClusterEngine(env); err != nil {
	 	message.WriteError(err)
	 	os.Exit(4)
	}
	// Call to initialize SampleEngine
  if err, sam = sample.NewSampleEngine(env); err != nil {
		message.WriteError(err)
		os.Exit(4)
  }
	// Call to initialize the CheckEngine
	if err, cks = check.NewCheckEngine(env); err != nil {
		message.WriteError(err)
		os.Exit(4)
	}
	// Call to initialize the ServiceEngine
	if err,srv = service.NewServiceEngine(env); err != nil {
		message.WriteError(err)
		os.Exit(4)
	}
	//
	//
	// Add the cluster engine's input channel as a service engine's output channel
	// That's the way how cluster engine will receive either node or services status
	if err := srv.AddOutputChannel(cltr.GetInputChannel(),"ClusterEngine"); err != nil {
		env.Output.WriteChWarn(err)
	}
	// Set the output sample channel for checks as the input's service one
	//cks.SetOutputSampleChan(srv.GetInputSampleChan())
	if err := cks.AddOutputChannel(srv.GetInputChannel(),"ServiceEngine"); err != nil {
		env.Output.WriteChWarn(err)
	}
	//cks.SetOutputSampleChan(srv.GetInputSampleChan())
	if err := cks.AddOutputChannel(sam.GetInputChannel(),"SampleEngine"); err != nil {
		env.Output.WriteChWarn(err)
	}

	switch(context.ExecutionMode){
	case "cluster":
		// prepare listen address for cluster node
		listenaddr := env.Context.Host+":"+strconv.Itoa(env.Context.Port)
		
		apisys := api.NewApiEngine(env)

		if err = cks.StartCheckEngine(nil); err != nil {
			env.Output.WriteChError(err)
			os.Exit(4)
		}

		webconsole := ui.NewUI(env, listenaddr)
		if err := cltr.AddOutputChannel(webconsole.GetInputChannel(),"UI"); err != nil {
			env.Output.WriteChError(err.Error())
		}
		webconsole.AddRoutes(apisys.GetRoutes())
		webconsole.StartUI()

		break
	case "standalone":
		message.Write("\t# That's Verdmell in standalone mode #\n")

		checkObj := new(check.CheckObject)
		//execute an isolated check
		if context.ExecuteCheck != "" {
			if err, checkObj = cks.GetCheckObjectByName(context.ExecuteCheck); err != nil {
				env.Output.WriteChError(err)
				os.Exit(4)	
			}	
			if err = cks.StartCheckEngine(checkObj); err != nil {
				env.Output.WriteChError(err)
				os.Exit(4)
			}
		// execute checks from group
		} else if  context.ExecuteCheckGroup != "" {
			var checks []string 
			if err, checks = cks.GetCheckgroupByName(context.ExecuteCheckGroup); err != nil {
				env.Output.WriteChError(err)
				os.Exit(4)	
			}
			if err = cks.StartCheckEngine(checks); err != nil {
				env.Output.WriteChError(err)
				os.Exit(4)
			}

		//execute all checks
		} else {
			if err = cks.StartCheckEngine(nil); err != nil {
				env.Output.WriteChError(err)
				os.Exit(4)
			}
		}

		// achieve required status
		if err, exitStatus = srv.GetServiceStatus(env.Context.Service); err != nil{
			env.Output.WriteChError(err)
			os.Exit(4)
		}

		_,hummanstatus := srv.GetServicesStatusHuman(env.Context.Service)
		message.Write(hummanstatus)

	}//end switch

	os.Exit(exitStatus)
}