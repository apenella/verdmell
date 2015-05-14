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
)
//
//Environment data type
//struct for config item into config json
type Environment struct{
		// Current node setup
		Setup setupObject
		// Output manager
		Output *message.Message
		// Context execution parameters
		context currentContext
}
//
// Methods for Environment
//

//
// Constructor for Environment
func NewEnvironment() (error, *Environment) {
	//object to dump configuration
	var env = new(Environment)
	var context = new(currentContext)
	var setup = new(setupObject)

	var err error

	// Set the output manager. This will control all the output during the system life
	output := message.GetInstance(context.Loglevel)
	env.SetOutput(output)

	if err, context = newcurrentContext(output); err != nil {return err, nil}
	// Set the context to the environment
	env.SetContext(context)

	if err, setup = newSetupObject(context.SetupFile, context.ConfigFolder, output); err != nil {return err, nil}
	env.SetSetup(setup)

	// Validate Environment
	if err := env.validateEnvironment(); err != nil {return err, nil}

	output.WriteChDebug("(Environment::NewEnvironment) Node environment ready...")
	return nil,env
}

//
// Setters
//
// Set setup for the Environment
func (e *Environment) SetSetup(s *setupObject) {
	e.Setup = *s
}
// Set output for the Environment
func (e *Environment) SetOutput(o *message.Message){
	e.Output = o
}
// Set the context for the Environment
func (e *Environment) SetContext(c *currentContext) {
	e.context = *c
}

//
// Getters
//
// Get the setupObject from envirionment
func (e *Environment) GetSetup() *setupObject{
		return &e.Setup
}
// Get output from environment
func (e *Environment) GetOutput() *message.Message{
	return e.Output
}
// Get context from environment
func (e *Environment) GetContext() *currentContext{
		return &e.context
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
func (e *Environment) String() string{
		str := "{"
		s := e.GetSetup()
		str += s.String()
		c := e.GetContext()
		str += c.String()
		str += "}"
		return str
}