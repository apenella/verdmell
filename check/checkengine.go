/*
Package check is used by verdmell to manage the monitoring checks defined by user
*/
package check

import (
	"errors"
	"os"
	"strconv"
	"time"

	"verdmell/environment"
	"verdmell/sample"
	"verdmell/utils"

	"github.com/apenella/messageOutput"
)

// DEFAULT_TIMEOUT defines a general porpuse timeout value
const DEFAULT_TIMEOUT int = 60

// CHECKS_KEY is the key used to look up the checks on json structures
const CHECKS_KEY string = "checks"

var env *environment.Environment
var logger *message.Message

// CheckEngine struct
// The struct for CheckEngine has all Check, Checkgroup and samples information
type CheckEngine struct {
	// Map to storage the checks
	// Cks *Checks `json:"checks"`
	Checks map[string]*Check `json:"checks"`
	// Service Channel
	subscriptions map[chan interface{}]string `json: "-"`
	// Folder where are found check files
	checksFolder string `json: "-"`
}

//
// NewCheckEngine is used as a constructor and returns a CheckEngine instance
func NewCheckEngine() *CheckEngine {

	logger := message.GetMessager()
	logger.Debug("(CheckEngine::NewCheckEngine) Create new engine instance")

	eng := &CheckEngine{
		Checks:        make(map[string]*Check),
		subscriptions: make(map[chan interface{}]string),
	}

	return eng
}

//
// Interface Engine requirements

// Init
func (eng *CheckEngine) Init() error {
	var err error
	var cks *Checks

	// initialize logger.Er
	if logger == nil {
		logger = message.New(message.INFO, nil, 0)
	}

	if eng.checksFolder != "" {
		// Get defined checks
		if err = eng.LoadChecks(""); err != nil {
			return errors.New("(CheckEngine::Init) " + err.Error())
		}
		// validate checks and set the checks into check system
		if err = cks.ValidateChecks(nil); err != nil {
			return errors.New("(CheckEngine::Init) " + err.Error())
		}
	}

	return nil
}

// Run
func (eng *CheckEngine) Run() error {
	return nil
}

func (eng *CheckEngine) Start() error                      { return nil }
func (eng *CheckEngine) Stop() error                       { return nil }
func (eng *CheckEngine) Status() int                       { return 0 }
func (eng *CheckEngine) GetID() uint                       { return uint(0) }
func (eng *CheckEngine) GetName() string                   { return "" }
func (eng *CheckEngine) GetDependencies() []uint           { return nil }
func (eng *CheckEngine) GetInputChannel() chan interface{} { return nil }
func (eng *CheckEngine) GetStatus() uint                   { return uint(0) }
func (eng *CheckEngine) SetStatus(s uint)                  {}

//
// Specific methods

//
// SayHi:
func (eng *CheckEngine) SayHi() {
	logger.Info("(CheckEngine::SayHi) Hi! I'm your new check engine instance")
}

//
// Subscribe:
func (eng *CheckEngine) Subscribe(o chan interface{}, desc string) error {
	logger.Debug("(CheckEngine::Subscribe) ", desc)

	channels := eng.subscriptions
	if _, exist := channels[o]; !exist {
		channels[o] = desc
	} else {
		return errors.New("(CheckEngine::Subscribe) You are trying to add an existing channel")
	}

	return nil
}

//
// LoadChecks
func (eng *CheckEngine) LoadChecks(folder string) error {
	checksQueue := make(chan []*Check)
	checksQueueErr := make(chan error)

	// files is an array with all files found inside the folder
	if folder == "" {
		folder = eng.checksFolder
	}
	files := utils.GetFolderFiles(folder)

	// call the goroutine for each file
	for _, file := range files {
		go func(f os.FileInfo) {
			checks, err := retrieveChecksFromFile(folder + string(os.PathSeparator) + f.Name())
			if err != nil {
				checksQueueErr <- err
			} else {
				checksQueue <- checks
			}
		}(file)
	}

	// analize response
	select {
	case err := <-checksQueueErr:
		logger.Info("(CheckEngine::LoadChecks) " + err.Error())
	case checks := <-checksQueue:
		for _, check := range checks {
			err := check.ValidateCheck()
			if err != nil {
				logger.Warn("(CheckEngine::LoadChecks) " + err.Error())
			} else {
				eng.AddCheck(check)
			}
		}
	case <-time.After(time.Duration(DEFAULT_TIMEOUT) * time.Second):

	}

	return nil
}

//
// retrieveChecksFromFile
func retrieveChecksFromFile(file string) ([]*Check, error) {
	var checks map[string][]*Check

	// extract the content from the file and dumps it on the CHecks object
	if err := utils.LoadJSONFile(file, &checks); err != nil {
		return nil, errors.New("(checkengine::retrieveChecksFromFile) Checks from '" + file + "' could not be retrieved. " + err.Error())
	}

	return checks[CHECKS_KEY], nil
}

//
// notify: function write samples to defined samples
func (eng *CheckEngine) notify(s *sample.CheckSample) {
	logger.Debug("(CheckEngine::notify)")
	for o, desc := range eng.subscriptions {
		logger.Debug("(CheckEngine::notify) [" + strconv.Itoa(int(s.GetTimestamp())) + "] Notify sample '" + s.GetCheck() + "' with exit '" + strconv.Itoa(s.GetExit()) + "' on channel '" + desc + "'")
		o <- s
	}
}

//
// AddCheck: method add a new check to the Checks struct
func (eng *CheckEngine) AddCheck(check *Check) error {
	eng.Checks[check.Name] = check

	return nil
}

//
// InitCheckRunningQueues: prepares each checkobject to be run
// func (eng *CheckEngine) InitCheckRunningQueues() error {
// 	logger.Debug("(CheckEngine::InitCheckRunningQueues)")
// 	cks := eng.Checks

// 	for _, obj := range cks.GetCheck() {
// 		go func(checkObj *CheckObject) {
// 			logger.Debug("(CheckEngine::InitCheckRunningQueues) CheckQueue for '" + checkObj.GetName() + "'")
// 			checkObj.StartQueue()
// 		}(obj)
// 	}
// 	return nil
// }

//
// Start: will determine which kind of check has been required by user and start the checks
// func (eng *CheckEngine) Start(i interface{}) error {
// 	logger.Debug("(CheckEngine::Start)")
// 	endChan := make(chan bool)
// 	defer close(endChan)
// 	errChan := make(chan error)
// 	defer close(errChan)

// 	// the next will be only used during ec or eg
// 	// check will contain the Check configurations
// 	check := new(Checks)
// 	// checks will contain all the CheckObject definition
// 	checks := make(map[string]*CheckObject)

// 	switch req := i.(type) {
// 	case *CheckObject:
// 		logger.Debug("(CheckEngine::Start) Starting the check '" + req.String() + "'")
// 		//add the check to be executed
// 		checks[req.GetName()] = req
// 		//add the check dependencies
// 		for _, dependency := range req.GetDepend() {
// 			if err, checkObj := eng.Cks.GetCheckObjectByName(dependency); err != nil {
// 				return err
// 			} else {
// 				if _, exist := checks[dependency]; !exist {
// 					checks[dependency] = checkObj
// 				}
// 			}
// 		}
// 		check.SetCheck(checks)
// 		// run a goroutine for each checkObject and write the result to the channel
// 		go func() {
// 			if err := check.StartCheckTaskPools(); err != nil {
// 				errChan <- err
// 			}
// 			endChan <- true
// 		}()

// 		select {
// 		case <-endChan:
// 			logger.Debug("(CheckEngine::Start) All Pools Finished")
// 		case err := <-errChan:
// 			return err
// 		}

// 	case []string:
// 		logger.Debug("(CheckEngine::Start) Running a Checkgroup")
// 		for _, checkname := range req {
// 			cks := eng.GetChecks()

// 			if err, checkObj := cks.GetCheckObjectByName(checkname); err != nil {
// 				return err
// 			} else {
// 				if _, exist := checks[checkname]; !exist {
// 					checks[checkname] = checkObj
// 				}
// 				//add the check dependencies
// 				for _, dependency := range checkObj.GetDepend() {
// 					if err, checkObjdependency := eng.Cks.GetCheckObjectByName(dependency); err != nil {
// 						return err
// 					} else {
// 						if _, exist := checks[dependency]; !exist {
// 							checks[dependency] = checkObjdependency
// 						}
// 					}
// 				}
// 			}
// 		}
// 		check.SetCheck(checks)
// 		// run a goroutine for each checkObject and write the result to the channel
// 		go func() {
// 			// startCheckTaskPools requiere the SAmple system to sent sample to it and OutputSampleChan to send samples to ServiceSystem
// 			if err := check.StartCheckTaskPools(); err != nil {
// 				errChan <- err
// 			}
// 			endChan <- true
// 		}()

// 		select {
// 		case <-endChan:
// 			logger.Debug("(CheckEngine::Start) All Pools Finished")
// 		case err := <-errChan:
// 			return err
// 		}

// 	default:
// 		checks := eng.GetChecks()
// 		if err := checks.StartCheckTaskPools(); err != nil {
// 			return err
// 		}
// 		logger.Debug("(CheckEngine::Start) All Pools Finished")
// 	}

// 	return nil
// }

//
// sendSamples: method that send samples to other engines
func (eng *CheckEngine) sendSample(s *sample.CheckSample) error {
	logger.Debug("(CheckEngine::sendSample)[" + strconv.Itoa(int(s.GetTimestamp())) + "] Send sample for '" + s.GetCheck() + "' check with exit '" + strconv.Itoa(s.GetExit()) + "'")
	sampleEngine := env.GetSampleEngine().(*sample.SampleEngine)

	// send samples to ServiceEngine
	// GetSample return an error if no samples has been add before for that check
	if err, sam := sampleEngine.GetSample(s.GetCheck()); err != nil {
		// sending sample to service using the output channel
		eng.notify(s)
	} else {
		// if a sample for that exist
		// the sample will not send to service system unless it has modified it exit status
		if sam.GetTimestamp() < s.GetTimestamp() {
			// sending sample to service using the output channel
			eng.notify(s)
		}
	}

	return nil
}

//
// ListCheckNames: returns an array with the check namess defined on Checks object
func (eng *CheckEngine) ListCheckNames() []string {
	checkNames := []string{}
	for name := range eng.Checks {
		checkNames = append(checkNames, name)
	}
	return checkNames
}

//
// IsDefined: return if a check is defined
func (eng *CheckEngine) IsDefined(name string) bool {
	_, ok := eng.Checks[name]
	return ok
}

//
// GetCheckObjectByName: returns a check object gived a name
func (eng *CheckEngine) GetCheckObjectByName(checkname string) (*Check, error) {
	return eng.Checks[checkname], nil
}

//
// GetAllChecks: return all checks
// func (eng *CheckEngine) GetAllChecks() ([]byte, error) {
// 	var checks *Checks

// 	if checks = eng.GetChecks(); checks == nil {
// 		msg := "(CheckEngine::GetAllChecks) There are no checks defined."
// 		logger.Debug(msg)
// 		return errors.New(msg), nil
// 	}

// 	return utils.ObjectToJsonByte(checks)
// }

//
// GetCheck: return a checks
// func (eng *CheckEngine) GetCheck(name string) ([]byte, error) {
// 	var checks *Checks
// 	var check map[string]*CheckObject
// 	var obj *CheckObject
// 	var exist bool

// 	// Get Checks attribute from CheckEngine
// 	if checks = eng.GetChecks(); checks == nil {
// 		msg := "(CheckEngine::GetCheck) There are no checks defined."
// 		logger.Debug(msg)
// 		return errors.New(msg), nil
// 	}
// 	// Get Check map from Checks
// 	if check = checks.GetCheck(); checks == nil {
// 		msg := "(CheckEngine::GetCheck) There are no checks defined."
// 		logger.Debug(msg)
// 		return errors.New(msg), nil
// 	}
// 	// Get CheckObject from the check's map
// 	if obj, exist = check[name]; !exist {
// 		msg := "(ServiceEngine::GetCheck) The check '" + name + "' is not defined."
// 		logger.Debug(msg)
// 		return errors.New(msg), nil
// 	}

// 	return utils.ObjectToJsonByte(obj)
// }

//
// Common methods

//
// String: convert a Checks object to string
func (eng *CheckEngine) String() string {
	var err error
	var str string

	if err, str = utils.ObjectToJsonString(eng); err != nil {
		return err.Error()
	}

	return str
}
