/*
 Package configuration manage all about configuration.
*/

package configuration

import (
	"errors"
	"os"

	"verdmell/utils"
)

// environment variable VERDMELL_HOME
const (
	VERDMELL_HOME string = "VERDMELL_HOME"
	CONF_DIR      string = "conf.d"
	CONF_FILE     string = "config.json"
	LOAD_TIMEOUT  int    = 60
)

//
// ChecksConfiguration has the data related to check configuration
type ChecksConfiguration struct {
	// Folder is where to locate checks
	Folder string `json:"folder"`
	// MinInterval is the minimum interval between check executions
	MinInterval int `json:"min_interval"`
	// LoadTimeout is the max time required to load checks from files
	LoadTimeout int `json:"load_timeout"`
}

//
// ServicesConfiguration has the data related to service configuration
type ServicesConfiguration struct {
	// folder to locate services
	Folder string `json:"folder"`
	// LoadTimeout is the max time required to load services from files
	LoadTimeout int `json:"load_timeout"`
}

//
// Configuration data
type Configuration struct {
	// Name for node
	Name string `json:"name"`
	// node's IP
	IP string `json:"ip"`
	// port
	Port int `json:"port"`
	//Even it is named Cluster, it refears to cluster nodes
	Cluster []string `json:"cluster"`
	// Configuration for checks
	Checks *ChecksConfiguration `json:"checks"`
	// Configuration for services
	Services *ServicesConfiguration `json:"services"`
}

//
// NewConfiguration constructs a configuration instance
func NewConfiguration(file string, dir string) (*Configuration, error) {
	var err error
	configuration := new(Configuration)

	// initialize global variables
	verdmellHome := os.Getenv(VERDMELL_HOME)
	if verdmellHome == "" {
		if verdmellHome, err = os.Getwd(); err != nil {
			verdmellHome = "."
		}
	}
	configurationDir := verdmellHome + string(os.PathSeparator) + CONF_DIR
	configurationFile := configurationDir + string(os.PathSeparator) + CONF_FILE

	// change configuration dir and file when a dir is not empty
	if dir != "" {
		configurationDir = dir
		configurationFile = configurationDir + string(os.PathSeparator) + CONF_FILE
	}

	// change configuration file when file is not empte
	if file != "" {
		configurationFile = configurationDir + string(os.PathSeparator) + file
	}

	// load configuration from file
	err = utils.LoadJSONFile(configurationFile, configuration)
	if err != nil {
		return nil, err
	}

	// Set path to check folder
	configuration.Checks.Folder = configurationDir + string(os.PathSeparator) + configuration.Checks.Folder
	// Set path to check folder
	configuration.Services.Folder = configurationDir + string(os.PathSeparator) + configuration.Services.Folder

	// Set timeout to load checks
	if configuration.Checks.LoadTimeout < 0 {
		configuration.Checks.LoadTimeout = LOAD_TIMEOUT
	}
	// Set timeout to load services
	if configuration.Services.LoadTimeout < 0 {
		configuration.Services.LoadTimeout = LOAD_TIMEOUT
	}

	err = configuration.ValidateConfiguration()
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

//
// ValidateConfiguration validates configuration object
func (c *Configuration) ValidateConfiguration() error {

	// validate name properties
	if c.Name == "" {
		return errors.New("(Configuration::ValidateConfiguration) Undefined name properties on configuration file")
	}
	if c.Port < 0 || c.Port > 65535 {
		return errors.New("(Configuration::ValidateConfiguration) Invalid port definition")
	}
	// validate checks properties
	if err := c.validateChecks(); err != nil {
		return errors.New("(Configuration::ValidateConfiguration) " + err.Error())
	}
	// validate services properties
	if err := c.validateServices(); err != nil {
		return errors.New("(Configuration::ValidateConfiguration) " + err.Error())
	}

	return nil
}

// validateChecks method to validate Checks information
func (c *Configuration) validateChecks() error {
	if c.Checks == nil {
		return errors.New("(Configuration::validateChecks) Undefined checks properties on configuration file")
	}
	if _, err := os.Stat(c.Checks.Folder); err != nil {
		return errors.New("(Configuration::validateChecks) Folder '" + c.Checks.Folder + "' does not exist")
	}

	if c.Checks.MinInterval < 0 {
		return errors.New("(Configuration::validateChecks) Invalid minimum interval")
	}

	if c.Checks.LoadTimeout < 0 {
		return errors.New("(Configuration::validateChecks) Invalid load timeout")
	}

	return nil
}

// validateServices method to validate Services information
func (c *Configuration) validateServices() error {
	if c.Services == nil {
		return errors.New("(Configuration::validateServices) Undefined services properties on configuration file")
	}
	if _, err := os.Stat(c.Services.Folder); err != nil {
		return errors.New("(Configuration::validateServices) Folder '" + c.Services.Folder + "' does not exist")
	}

	if c.Services.LoadTimeout < 0 {
		return errors.New("(Configuration::validateServices) Invalid load timeout")
	}

	return nil
}

//
// Common methods

// String method transform the Configuration to string
func (c *Configuration) String() string {
	var err error
	var str string

	str, err = utils.ObjectToJSONString(c)
	if err != nil {
		return err.Error()
	}

	return str
}
