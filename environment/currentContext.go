/*
Environment: manage all data related with the execution and any thing around it.

-Environment
-SetupObject
-currentContext
*/
package environment

import (
	"os"
	"errors" 
	"flag"
	"strconv"
	"github.com/apenella/messageOutput"
	"verdmell/utils"
)

//struc for running Environment
type currentContext struct{
	// configuration folder
	ConfigFolder string
	// configuration file
	SetupFile string
	// service to get information for
	Service string
	// execute the indicated check
	ExecuteCheck string
	// execute the indicated checkgroup
	ExecuteCheckGroup string
	// loglevel definition
	/*
		0: information
		1: warn
		2: error
		3: debug
	*/
	Loglevel int  
	// execution mode
	/*
		standalone
		cluster
	*/
	ExecutionMode string

	// nodes that belongs to cluster
	Nodes []string
	// host to anchor to server mode
	Host string
	// port to anchor to server mode
	Port int
	// output manager
	output *message.Message
}

func newCurrentContext(output *message.Message) (error, *currentContext) {
	context := new(currentContext)

	var loglevel int
	var configFolder string
	var setupFile string
	var service string
	var executeCheck string
	var executeCheckGroup string
	var executionMode string
	var port int
	var host string

	flag.IntVar(&loglevel,"l",0,"Loglevel definition\n\t0 - info\n\t1 - warn\n\t2 - error\n\t3 - debug.")  
	flag.StringVar(&configFolder,"d","./conf.d","Root configuration folder.")
	flag.StringVar(&setupFile,"c","config.json","Configuration file.")
	flag.StringVar(&service,"s","","Specific service which you ask its status for")
	flag.StringVar(&executeCheck,"ec","","Execute the indicated check.")
	flag.StringVar(&executeCheckGroup,"eg","","Execute the indicated check group.")
	flag.StringVar(&executionMode,"m","standalone","Execution mode indicates how to run verdmell.\n\t-standalone: return the health status ondemand\n\t-cluster: start a service which is listening for health status requests")
	flag.IntVar(&port,"p",5497,"Set a custom port for the cluster mode")
	flag.StringVar(&host,"h","0.0.0.0","Set a custom IP for the server mode")
	flag.Parse()

	output.SetLogLevel(loglevel)

	context = &currentContext{
		ConfigFolder: configFolder,
		SetupFile: setupFile,
		Service: service,
		ExecuteCheck: executeCheck,
		ExecuteCheckGroup: executeCheckGroup,
		Loglevel: loglevel,
		ExecutionMode: executionMode,
		Host: host,
		Port: port,
		output: output,
	}

	err := context.validatecurrentContext()

	return err, context
}

//
// Specific methods
// validatecurrentContext check each context parameter to ensure its correctness
func (c *currentContext) validatecurrentContext() error {
	c.output.WriteChDebug("(currentContext::validatecurrentContext) validation current context")
	c.output.WriteChDebug("(currentContext::validatecurrentContext) configFolder: "+c.ConfigFolder)
	c.output.WriteChDebug("(currentContext::validatecurrentContext) configFile: "+c.SetupFile)
	c.output.WriteChDebug("(currentContext::validatecurrentContext) execution mode: "+c.ExecutionMode)
	c.output.WriteChDebug("(currentContext::validatecurrentContext) service: "+c.Service)
	c.output.WriteChDebug("(currentContext::validatecurrentContext) execute check: "+c.ExecuteCheck)
	c.output.WriteChDebug("(currentContext::validatecurrentContext) execute checkgroup: "+c.ExecuteCheckGroup)
	c.output.WriteChDebug("(currentContext::validatecurrentContext) execute IP: "+c.Host)
	c.output.WriteChDebug("(currentContext::validatecurrentContext) execute port: "+strconv.Itoa(c.Port))

	// configuration folder
	if err := utils.FileExist(c.ConfigFolder); os.IsNotExist(err) {return err}
	// configuration file
	if err := utils.FileExist(c.ConfigFolder+string(os.PathSeparator)+c.SetupFile); os.IsNotExist(err) {return err}
	// execute the indicated check
	// at this point we couldn't validate if the check exists because chekcs haven't already been loaded
	if c.ExecuteCheck != "" && c.ExecuteCheckGroup != "" {return errors.New("You should decide between execute a check or a checkgroup")}  
	// loglevel definition
	if c.Loglevel < 0 || c.Loglevel > 3 {return errors.New("Undefined loglevel mode")}
	// execution mode
	if c.ExecutionMode != "standalone" && c.ExecutionMode != "cluster" { return errors.New("The execution mode chosen is unknown by Verdmell")}
	// nodes only take sense in cluster mode
	if c.ExecutionMode == "standalone" && len(c.Nodes) > 0 {return errors.New("The node only have sense on cluster's execution mode")}
	// host to anchor to server mode
	if err := utils.IsLocalIPAddress(c.Host); err != nil {return err}
	// port to anchor to server mode
	if c.Port < 0 || c.Port > 65535 {return errors.New("Is not possible to use a port out of 0..65535 range")}
	//if c.ExecutionMode != "server" && (c.Port != "" || c.Host != "") { return errors.New("There is no sense using the host or port flag for non server mode")}
	return nil
}

//
// Common methods

// method to transform the currentContext to string
func (c *currentContext) String() string{
	if err, str := utils.ObjectToJsonString(c); err != nil{
		return err.Error()
	} else{
		return str
	}
}

//#####################################################################################################
