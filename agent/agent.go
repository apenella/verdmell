package agent

import (
	"verdmell/configuration"
	"verdmell/engine"
	"verdmell/environment"
	"verdmell/utils"

	"github.com/apenella/messageOutput"
)

var env *environment.Environment

/*
	Agent is the element that coordinates all components
*/
type Agent struct{
	// loglevel for agent
	Loglevel int `json: "loglevel"`
	// configuration file name
	Configfile string `json: "configuration_file"`
	// folder to place configuration
	Configdir string `json: "configuration_dir"`

	// List of engines
	Engines map[uint]engine.Engine `json: "-"`
}

//
// Common methods
//---------------------------------------------------------------------

/*
	Start
*/
func (a* Agent) Start() error {
	
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

	env.Output.WriteChInfo("Agent Start")
	return nil
}

/*
	Status
*/
func (a* Agent) Status() error {
	env.Output.WriteChInfo("Agent Status")
	return nil
}

// String method transform the Configuration to string
func (a *Agent) String() string {
	if err, str := utils.ObjectToJsonString(a); err != nil{
		return err.Error()
	} else{
		return str
	} 
}