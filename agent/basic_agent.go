package agent

import (
	"errors"

	"verdmell/check"
	"verdmell/configuration"
	"verdmell/environment"
	"verdmell/sample"
	"verdmell/service"
	"verdmell/utils"

	"github.com/apenella/messageOutput"
)

var env *environment.Environment

type BasicAgent struct {
	Agent

	// contain list of checks to be run
	Checks StringList `json: "checks"`
}

func (a *BasicAgent) Start() error {
	var err error
	env = &environment.Environment {
		Output: message.GetInstance(a.Loglevel),
	}

	// generate a configuration
	if err, configuration := configuration.NewConfiguration(a.Configfile, a.Configdir, env.Output); err != nil {
		env.Output.WriteChError(err)
		return err
	} else {
		env.SetConfig(configuration)
	}

	// Call to initialize SampleEngine
  	if err, a.Sam = sample.NewSampleEngine(env); err != nil {
		env.Output.WriteChError(err)
		return err
  	}
	// Call to initialize the CheckEngine
	if err, a.Cks = check.NewCheckEngine(env); err != nil {
		env.Output.WriteChError(err)
		return err
	}
	// Call to initialize the ServiceEngine
	if err, a.Srv = service.NewServiceEngine(env); err != nil {
		message.WriteError(err)
		return err
	}

	//subscribre samples to checks

	if err := a.Cks.Subscribe(a.Sam.GetInputChannel(),"SampleEngine"); err != nil {
		env.Output.WriteChError(err)
		return err
	}
	//subscribre services to checks
	if err := a.Cks.Subscribe(a.Srv.GetInputChannel(),"ServiceEngine"); err != nil {
		env.Output.WriteChError(err)
		return err
	}

	if err = a.Cks.Start(a.Checks); err != nil {
		env.Output.WriteChError(err)
		return err
	}

	// for _,c := range cks.GetChecks().ListCheckNames() {
	// 	env.Output.WriteChInfo(sam.GetSample(c))
	// }

	return nil
}

// Validate state of agent
func (a *BasicAgent) Validate() error {
	if a.Cks == nil {
		return errors.New("(BasicAgent::Status) Null CheckEngine")	
	}
	if a.Sam == nil {
		return errors.New("(BasicAgent::Status) Null SampleEngine")	
	}
	if a.Srv == nil {
		return errors.New("(BasicAgent::Status) Null ServiceEngine")	
	}

	return nil
}

// return status
func (a *BasicAgent) Status() error {

	if err := a.Validate(); err != nil {
		return err
	}
	
	if err, status := a.Srv.GetServicesStatusHuman(env.Config.Name); err != nil {
		env.Output.WriteChError(err)
		return err
	} else {
		message.Write(status)		
	}
	return nil
}

//
// Common methods
//---------------------------------------------------------------------

// String method transform the Configuration to string
func (a *BasicAgent) String() string {
	if err, str := utils.ObjectToJsonString(a); err != nil{
		return err.Error()
	} else{
		return str
	} 
}