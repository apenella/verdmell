/*
Environment: manage all data related with the execution and any thing around it.

-Environment
-SetupObject
-currentContext
*/
package environment

import (
	"errors"
	"github.com/apenella/messageOutput"
	"verdmell/utils"
)
//
//Environment data type
//struct for config item into config json
type Environment struct{
	// Current node setup
	Setup *setupObject
	// Output manager
	Output *message.Message
	// Context execution parameters
	Context *currentContext

	//ChecksEngine
	CheckEngine interface{} `json:"checks"`
	//SampleEngine
	SampleEngine interface{} `json:"samples"`
	//ServiceEngine
	ServiceEngine interface{} `json:"services"`
	//ClusterEngine
	ClusterEngine interface{} `json:"cluster"`
}
//
// Methods for Environment
//

//
// Constructor for Environment
func NewEnvironment() (error, *Environment) {
	var err error
	var context = new(currentContext)
	var setup = new(setupObject)

	output := message.GetInstance(context.Loglevel)
	if output == nil {
		return errors.New("(Environment::NewEnvironment) The OutputMessage instance is null"),nil
	}

	if err, context = newCurrentContext(output); err != nil {return err, nil}
	if err, setup = newSetupObject(context.SetupFile, context.ConfigFolder, output); err != nil {return err, nil}

	env := &Environment{
		Setup: setup,
		Context: context,
		Output: output,
	}
	
	if env.Context.Service == "" {
		env.Context.Service = env.Setup.Hostname
	}

	// Validate Environment
	if err := env.validateEnvironment(); err != nil {return err, nil}

	output.WriteChDebug("(Environment::NewEnvironment) Node environment ready...")
	return nil,env
}

//#
//# Getters and Setters
//#----------------------------------------------------------------------------------------

// Set setup for the Environment
func (e *Environment) SetSetup(s *setupObject) {
	e.Output.WriteChDebug("(Environment::SetSetup)")
	e.Setup = s
}
// Set output for the Environment
func (e *Environment) SetOutput(o *message.Message){
	e.Output.WriteChDebug("(Environment::SetOutput)")
	e.Output = o
}
// Set the context for the Environment
func (e *Environment) SetContext(c *currentContext) {
	e.Output.WriteChDebug("(Environment::SetContext)")
	e.Context = c
}
// Set the CheckEngine for the Environment
func (e *Environment) SetCheckEngine(c interface{}) {
	e.Output.WriteChDebug("(Environment::SetCheckEngine)")
	e.CheckEngine = c
}
// Set the SampleEngine for the Environment
func (e *Environment) SetSampleEngine(s interface{}) {
	e.Output.WriteChDebug("(Environment::SetSampleEngine)")
	e.SampleEngine = s
}
// Set the ServiceEngine for the Environment
func (e *Environment) SetServiceEngine(s interface{}) {
	e.Output.WriteChDebug("(Environment::SetServiceEngine)")
	e.ServiceEngine = s
}
// Set the ServiceEngine for the Environment
func (e *Environment) SetClusterEngine(s interface{}) {
	e.Output.WriteChDebug("(Environment::SetClusterEngine)")
	e.ClusterEngine = s
}

// Get the setupObject from envirionment
func (e *Environment) GetSetup() *setupObject{
	e.Output.WriteChDebug("(Environment::GetSetup)")
	return e.Setup
}
// Get output from environment
func (e *Environment) GetOutput() *message.Message{
	e.Output.WriteChDebug("(Environment::GetOutput)")
	return e.Output
}
// Get context from environment
func (e *Environment) GetContext() *currentContext{
	e.Output.WriteChDebug("(Environment::GetContext)")
	return e.Context
}
// Get CheckEngine from environment
func (e *Environment) GetCheckEngine() interface{} {
	e.Output.WriteChDebug("(Environment::GetCheckEngine)")
	return e.CheckEngine
}
// Get SampleEngine from environment
func (e *Environment) GetSampleEngine() interface{} {
	e.Output.WriteChDebug("(Environment::GetSampleEngine)")
	return e.SampleEngine
}
// Get ServiceEngine from environment
func (e *Environment) GetServiceEngine() interface{} {
	e.Output.WriteChDebug("(Environment::GetServiceEngine)")
	return e.ServiceEngine
}
// Get ClusterEngine from environment
func (e *Environment) GetClusterEngine() interface{} {
	e.Output.WriteChDebug("(Environment::GetClusterEngine)")
	return e.ClusterEngine
}

//
// Specific methods
//---------------------------------------------------------------------

//
//# validateEnvironment: method to validate configuration objecte
func (e *Environment) validateEnvironment() error {
		if e == nil {
				err := errors.New("(Environment::validateEnvironment) Configuration object is empty due to error on configuration file")
				return err
		}

		s := e.GetSetup()
		if err := s.validateSetupObject(); err == nil {return err}

		return nil
}
//
//# GetNodeInfo from node
func (e *Environment) GetNodeInfo() (error,[]byte) {
	e.Output.WriteChDebug("(Environment::GetNodeInfo)")

	environment := make(map[string]interface{})

	environment["checks"] = e.GetCheckEngine()
	environment["samples"] = e.GetSampleEngine()
	environment["services"] = e.GetServiceEngine()

	return nil,utils.ObjectToJsonByte(environment)
	
}
//
//# GetCluster return all cluster nodes
func (e *Environment) GetCluster() []byte{
	return utils.ObjectToJsonByte(e.Setup.Cluster)
}
//
// Common methods
//---------------------------------------------------------------------

//
// String method to return the string
func (e *Environment) String() string {
	if err, str := utils.ObjectToJsonString(e); err != nil{
		return err.Error()
	} else{
		return str
	}
}

//####################################################################################################