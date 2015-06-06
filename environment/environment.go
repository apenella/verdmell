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
	//Check names
	Checks []string
	//Services names
	Services []string
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
	if err, context = newcurrentContext(output); err != nil {return err, nil}
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
	e.Setup = s
}
// Set output for the Environment
func (e *Environment) SetOutput(o *message.Message){
	e.Output = o
}
// Set the context for the Environment
func (e *Environment) SetContext(c *currentContext) {
	e.Context = c
}
// Set the Checks or the Environment
func (e *Environment) SetChecks(c []string) {
	e.Checks = c
}
// Set the Services for the Environment
func (e *Environment) SetServices(s []string) {
	e.Services = s
}

// Get the setupObject from envirionment
func (e *Environment) GetSetup() *setupObject{
		return e.Setup
}
// Get output from environment
func (e *Environment) GetOutput() *message.Message{
	return e.Output
}
// Get context from environment
func (e *Environment) GetContext() *currentContext{
		return e.Context
}
// Get Checks from environment
func (e *Environment) GetChecks() []string{
		return e.Checks
}
// Get Services from environment
func (e *Environment) GetServices() []string{
		return e.Services
}

//
// Specific methods
//---------------------------------------------------------------------

//
// method to validate configuration objecte
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
// Common methods
//---------------------------------------------------------------------

//
// String method to return the string
func (e *Environment) String() string {
	return utils.ObjectToJsonString(e)
}

//####################################################################################################