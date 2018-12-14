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
	"os"
	"verdmell/utils"
)

//setupObject data type
//struct for setupObject contain the system data required during execution like, folders, hostname or ip.
type setupObject struct {
	Checksfolder   string `json:"checksfolder"`
	Servicesfolder string `json: "servicesfolder"`
	Hostname       string `json:"hostname"`
	Ip             string `json:"ip"`
	//Even it is named Cluster, it refears to cluster nodes
	Cluster map[string]string `json:"cluster"`
	// output manager
	output *message.Message
}

func newSetupObject(file string, folder string, output *message.Message) (error, *setupObject) {
	setup := new(setupObject)
	setup.output = output

	// Set path to setupFile
	file = folder + string(os.PathSeparator) + file
	output.WriteChDebug("(SetupObject::newSetupObject) Configuratin file '" + file + "'")
	// Dump setup data to Environment
	if err := utils.LoadJSONFile(file, setup); err != nil {
		return errors.New("(SetupObject::newSetupObject) " + err.Error()), nil
	}

	// Set path to check folder
	setup.Checksfolder = folder + string(os.PathSeparator) + setup.Checksfolder
	output.WriteChDebug("(SetupObject::newSetupObject) Checks' definition on '" + setup.Checksfolder + "'")
	// Set path to check folder
	setup.Servicesfolder = folder + string(os.PathSeparator) + setup.Servicesfolder
	output.WriteChDebug("(SetupObject::newSetupObject) Services' definition on '" + setup.Servicesfolder + "'")

	output.WriteChDebug("(SetupObject::newSetupObject) Created'")
	return nil, setup
}

//
// Specific methods
//---------------------------------------------------------------------

// validate setup object content
func (s *setupObject) validateSetupObject() error {
	if err := s.validateChecksfolder(); os.IsNotExist(err) {
		return err
	}
	if err := s.validateServicesfolder(); os.IsNotExist(err) {
		return err
	}
	if err := s.validateHostInfo(); os.IsNotExist(err) {
		return err
	}

	return nil
}

// method to validate Host information
func (s *setupObject) validateChecksfolder() error {
	if _, err := os.Stat(s.Checksfolder); err != nil {
		err := errors.New("(setupObject::validateChecksfolder) Folder " + s.Checksfolder + " does not exist.")
		return err

	}
	s.output.WriteChDebug("(setupObject::validateChecksfolder) '" + s.Checksfolder + "'")
	return nil
}

// method to validate Host information
func (s *setupObject) validateServicesfolder() error {
	if _, err := os.Stat(s.Servicesfolder); err != nil {
		err := errors.New("(setupObject::validateServicesfolder) Folder " + s.Servicesfolder + " does not exist.")
		return err

	}
	s.output.WriteChDebug("(setupObject::validateServicesfolder) '" + s.Servicesfolder + "'")
	return nil
}

// method to validate Host information
func (s *setupObject) validateHostInfo() error {
	s.output.WriteChDebug("(setupObject::validateHostInfo) '" + s.Hostname + "'")
	//TODO
	return nil
}

func (s *setupObject) AddNodeToCluster(node string, url string) bool {
	s.output.WriteChDebug("(setupObject::AddNodeToCluster) Adding node to cluster '" + node + "' with url " + url)
	var exist bool

	if _, exist = s.Cluster[node]; !exist {
		s.Cluster[node] = url
	} else {
		s.output.WriteChDebug("(setupObject::AddNodeToCluster) The node '" + node + "' has already exist into the cluster")
	}

	return exist
}

//
// Common methods
//---------------------------------------------------------------------

// String method transform the setupObject to string
func (s *setupObject) String() string {
	var str string
	var err error

	str, err = utils.ObjectToJSONString(s)
	if err != nil {
		return err.Error()
	}

	return str
}
