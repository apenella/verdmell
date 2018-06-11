//
// Configuration: manage all about configuration.
//
// -Configuration

package configuration

import (
	"errors"
	"os"

	"verdmell/utils"
	"github.com/apenella/messageOutput"
)

// environment variable VERDMELL_HOME
const (
	VERDMELL_HOME string = "VERDMELL_HOME"
	CONF_DIR string = "conf.d"
	CONF_FILE string = "config.json"
)

var verdmell_home string
var configuration_dir string
var configuration_file string

// data related to check configuration
type ChecksConfiguration struct {
	// folder to locate checks
	Folder string `json: "folder"`
	// minimum interval
	MinInterval int `json: "min_interval"`
}

// data related to service configuration
type ServicesConfiguration struct {
	// folder to locate services
	Folder string `json: "folder"`
}

// data related to configuration
type Configuration struct {
	// Name for node
	Name string `json:"name"`
	// node's IP
	Ip string `json:"ip"`
	// port
	Port int `json:"port"`
	//Even it is named Cluster, it refears to cluster nodes
	Cluster []string `json:"cluster"`
	// Configuration for checks
	Checks  *ChecksConfiguration `json:"checks"`
	// Configuration for services
	Services *ServicesConfiguration `json: "services"`
	// output manager
	log *message.Message
}

func init() {
	var err error
	verdmell_home = os.Getenv(VERDMELL_HOME)
	if verdmell_home == "" {
		if verdmell_home, err = os.Getwd(); err != nil {
			verdmell_home = "."
		}
	}
	configuration_dir = verdmell_home+string(os.PathSeparator)+CONF_DIR
	configuration_file = configuration_dir+string(os.PathSeparator)+CONF_FILE
}

// NewConfiguration: create new instance for configuration
func NewConfiguration(file string, dir string, log *message.Message) (error, *Configuration){
	configuration := new(Configuration)
	configuration.log = log

	// change configuration dir and file when a dir is not empty
	if dir != "" {
		configuration_dir = dir
		configuration_file = configuration_dir+string(os.PathSeparator)+CONF_FILE
	}

	// change configuration file when file is not empte
	if file != "" {
		configuration_file = configuration_dir+string(os.PathSeparator)+file
	}

	// Dump setup data to Environment
	if err := utils.LoadJSONFile(configuration_file, configuration); err != nil {
		return err, nil
	}

	// Set path to check folder
	configuration.Checks.Folder = configuration_dir+string(os.PathSeparator)+configuration.Checks.Folder
	// Set path to check folder
	configuration.Services.Folder = configuration_dir+string(os.PathSeparator)+configuration.Services.Folder

	if err := configuration.ValidateConfiguration(); err != nil {
		return err, nil
	}

	return nil, configuration
}

//
// Getters
//
// func (c *Configuration) Name() string {
// 	return c.Name
// }
// func (c *Configuration) Ip() string {
// 	return c.Ip
// }
// func (c *Configuration) Port() int {
// 	return c.Port
// }
// func  (c *Configuration) Cluster() []string {
// 	return c.Cluster
// }
// func  (c *Configuration) Checks() *ChecksConfiguration {
// 	return c.Checks
// }
// func  (c *Configuration) Services() *ServicesConfiguration {
// 	return c.Services
// }

// ValidateConfiguration: validates configuration object
func (c *Configuration) ValidateConfiguration() error {

	// validate name properties
	if c.Name == "" {
		return errors.New("(Configuration::ValidateConfiguration) Undefined name properties on configuration file")
	}
	// validate checks properties
	if err := c.validateChecks(); err != nil {
		return err
	}
	// validate services properties
	if err := c.validateServices(); err != nil {
		return err
	}

	return nil
}

// validateChecks: method to validate Checks information
func (c *Configuration) validateChecks() error{
	if c.Checks == nil {
		return errors.New("(Configuration::validateChecks) Undefined checks properties on configuration file")
	}
	if _, err := os.Stat(c.Checks.Folder); err != nil {
		return errors.New("(Configuration::validateChecks) Folder '"+c.Checks.Folder+"' does not exist")
	}

	return nil
}
// validateServices: method to validate Services information
func (c *Configuration) validateServices() error{
	if c.Services == nil {
		return errors.New("(Configuration::validateServices) Undefined services properties on configuration file")
	}
	if _, err := os.Stat(c.Services.Folder); err != nil {
		return errors.New("(Configuration::validateServices) Folder '"+c.Services.Folder+"' does not exist")
	}

	return nil
}

//
// Common methods
//---------------------------------------------------------------------

// String method transform the Configuration to string
func (c *Configuration) String() string{
	if err, str := utils.ObjectToJsonString(c); err != nil{
		return err.Error()
	} else{
		return str
	}
}
