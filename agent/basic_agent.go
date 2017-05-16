package agent

import (
//	"errors"

	"verdmell/check"
	"verdmell/configuration"
	"verdmell/environment"
	"verdmell/sample"
	"verdmell/utils"

	"github.com/apenella/messageOutput"
)

type BasicAgent struct {
	Agent

	// contain list of checks to be run
	Checks StringList `json: "checks"`
}

func (a *BasicAgent) Start() error {
	var err error
	var sam *sample.SampleEngine
	var cks *check.CheckEngine

	env := &environment.Environment {
		Output: message.GetInstance(a.Loglevel),
	}

	if err, configuration := configuration.NewConfiguration(a.Configfile, a.Configdir, env.Output); err != nil {
		env.Output.WriteChError(err)
		return err
	} else {
		env.SetConfig(configuration)
	}

	// Call to initialize SampleEngine
  	if err, sam = sample.NewSampleEngine(env); err != nil {
		env.Output.WriteChError(err)
		return err
  	}
	// Call to initialize the CheckEngine
	if err, cks = check.NewCheckEngine(env); err != nil {
		env.Output.WriteChError(err)
		return err
	}

	//subscribre samples to checks
	if err := cks.Subscribe(sam.GetInputChannel(),"SampleEngine"); err != nil {
		env.Output.WriteChError(err)
	}

	if err = cks.Start(a.Checks); err != nil {
		env.Output.WriteChError(err)
		return err
	}

	for _,c := range cks.GetChecks().ListCheckNames() {
		env.Output.WriteChInfo(sam.GetSample(c))
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