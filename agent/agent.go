package agent

import (
	"verdmell/check"
	"verdmell/sample"
	"verdmell/service"
	"verdmell/utils"
)

//
type Agenter interface {
	Start() error
}

//
type Agent struct{
	// loglevel for agent
	Loglevel int `json: "loglevel"`
	// configuration file name
	Configfile string `json: "configuration_file"`
	// folder to place configuration
	Configdir string `json: "configuration_dir"`

	// Check Engine manages checks
	Cks *check.CheckEngine `json: "-"`
	// Sample Engine manages samples
	Sam *sample.SampleEngine `json: "-"`
	// Service Engine manages services
	Srv *service.ServiceEngine `json: "-"`
}

//
// Common methods
//---------------------------------------------------------------------

// String method transform the Configuration to string
func (a *Agent) String() string {
	if err, str := utils.ObjectToJsonString(a); err != nil{
		return err.Error()
	} else{
		return str
	} 
}