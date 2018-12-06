package agent

import (

	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"verdmell/configuration"
	"verdmell/engine"
	"verdmell/environment"
	"verdmell/utils"

	"github.com/apenella/messageOutput"
	"github.com/ryanuber/columnize"
)

// define agent exit status
const (
	OK = iota
	WARN
	ERROR
	UNKNOWN
)

const DEFAULT_TIMEOUT int = 60
const DEFAULT_INIT_TIMEOUT int = 60
const DEFAULT_READY_TIMEOUT int = 60
const DEFAULT_STOP_TIMEOUT int = 60

var env *environment.Environment

// Agent is the element that coordinates all components
type Agent struct{
	// environment
	Ui io.Writer `json: "-"`

	// loglevel for agent
	Loglevel int `json: "loglevel"`
	// configuration file name
	Configfile string `json: "configuration_file"`
	// folder to place configuration
	Configdir string `json: "configuration_dir"`

	// List of engines
	Engines map[uint]engine.Engine `json: "engines"`
	EngineStatus map[uint]uint `json: "-"`

	// Is the engines order to be run, decided by user
	RunOrder []uint `json: "-"`
	// contains the order which the engines will run
	RunVector []uint `json: "runvector"`

	// Timeout to waiting for engines initialization
	InitTimeout int `json: "inittimeout"`
	// Timeout to waiting for engines to be ready
	ReadyTimeout int `json: "readytimeout"`
	// Timeout to waiting for engines to be stopped
	StopTimeout int `json: "stoptimeout"`

	// channels to notify init
	initCh chan uint `json: "-"`
	initErrCh chan error `json: "-"`

	// channels to notify readiness
	readyCh chan uint `json: "-"`
	readyErrCh chan error `json: "-"`

	// channels to notify stop
	stopCh chan uint `json: "-"`
	stopErrCh chan error `json: "-"`
}

//
// Common methods
//

// Start
func (a *Agent) Start() (int, error) {
	var err error

	// initialize agent
	err = a.init()
	if err != nil {
		env.Output.WriteChError("(Agent::Start)",err)
		return ERROR, err
	}

	// validate whether there is a dependency loop on engines
	if err = a.validateGraphEngine(); err != nil {
		env.Output.WriteChError("(Agent::Start)",err)
		return ERROR, err
	}

	// Initialize all engines
	if err = a.initialize(); err != nil {
		// stop the agent before return
		a.Stop()
		return ERROR, err
	}

	// Set up subscriptions among all engines
	if err = a.setEnginesSubscriptions(); err != nil {
		// stop the agent before return
		a.Stop()
		return ERROR, err
	}

	// Run engines
	if err = a.run(); err != nil {
		// stop the agent before return
		a.Stop()
		return ERROR, err
	}

	return OK, nil
}

func (a *Agent) Stop() int {
	var err error

	// Initialize Engines
	for id, e := range a.Engines {
		e.SetStatus(engine.STOPPING)

		// initialize current engine
		go func(id uint, e engine.Engine) {
			err := e.Stop()
			// write on error init channel whether there is an error on Init
			if err != nil {
				env.Output.WriteChInfo("(Agent::Stop)",err)
				a.stopErrCh <- err
			}
			// write engine id on init channel to notify that engine is already initialized
			a.stopCh <- id
		}(id,e)
	}
	// Waiting for all engines to be stopped
	// define a timeout to wait for all engines initialization
	timeout := time.After(time.Duration(a.StopTimeout) * time.Second)
	for i:=0; i < len(a.Engines); i++ {
		select{
			// wait to receive an id to set as initialized the engine
			case id := <-a.stopCh:
				e := a.Engines[id]
				e.SetStatus(engine.STOPPED)

			// wait to receive an error
			case err = <-a.stopErrCh:
				env.Output.WriteChError("(Agent::Stop)",err)

			// define new error when timeout is reached an not all engines are initialized
			case <- timeout:
				msg := "(Agent::Stop) Not all engines have been stopped after "+strconv.Itoa(int(a.StopTimeout))+" seconds."
				err = errors.New(msg)
				env.Output.WriteChError(err)
			}
	}

	defer close(a.initCh)
	defer close(a.initErrCh)
	defer close(a.readyCh)
	defer close(a.readyErrCh)
	defer close(a.stopCh)
	defer close(a.stopErrCh)

	// return when an error is detected stopping engines
	if err != nil {
		return ERROR
	}

	//a.Status()
	return OK
}

// Status
func (a *Agent) Status() int {

	lines := []string{
		"ENGINE |  STATUS",
		//"------ | ------",
	}

	for _, e := range a.Engines {
		lines = append(lines, e.GetName()+" | "+engine.ToHummanStatus(e.GetStatus()))
	}

	fmt.Println(columnize.SimpleFormat(lines))

	return 0
}

// String method transform the Configuration to string
func (a *Agent) String() string {
	if err, str := utils.ObjectToJsonString(a); err != nil{
		return err.Error()
	} else{
		return str
	}
}

//
// Private methods

// Initialize data structures required to agent to be run
func (a *Agent) init() error {
	// initialize environment
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

	// initialize engine status data structure
	a.EngineStatus = make(map[uint] uint)

	// initialize run order structure
	a.RunOrder = []uint{}

	// initialize channel to receive just initialited engines id
	a.initCh = make(chan uint)
	// initialize channel to receive error detected during initialization
	a.initErrCh = make(chan error)

	// initialize channel to receive ready engines id
	a.readyCh = make(chan uint)
	// initialize channel to receive error when run
	a.readyErrCh = make(chan error)

	// initialize channel to receive just stopped engines id
	a.stopCh = make(chan uint)
	// initialize channel to receive error when stopping
	a.stopErrCh = make(chan error)

	// set initialization timeout
	if a.InitTimeout <= 0 {
		a.InitTimeout = DEFAULT_INIT_TIMEOUT
	}
	// set ready timeout
	if a.ReadyTimeout <= 0 {
		a.ReadyTimeout = DEFAULT_READY_TIMEOUT
	}

	// set stop timeout
	if a.StopTimeout <= 0 {
		a.StopTimeout = DEFAULT_STOP_TIMEOUT
	}

	return nil
}

/*
	Improved run method

	1 - Initialize engine
	2 - In case of dependency loop, return
	3 - If no reverse dependency, start engine
	4 - Else,
		4.1 Decrease reverse dependency counter
		5.2 If reverse dependency counter is 0, start reverse dependency
			5.2.1 when reverse dependency is ready, try to start dependent engines from 3.1
*/

// method to initialize engines
func (a *Agent) initialize() error {
	var err error

	// Initialize Engines
	for id, e := range a.Engines {
		//e.SetStatus(engine.INITIALIZING)

		// initialize current engine
		go func(id uint, e engine.Engine) {
			err := e.Init()
			// write on error init channel whether there is an error on Init
			if err != nil {
				env.Output.WriteChInfo("(Agent::initialize)",err)
				a.initErrCh <- err
			}
			// write engine id on init channel to notify that engine is already initialized
			a.initCh <- id
		}(id,e)
	}

	// Waiting for all engines to be initialized
	// define a timeout to wait for all engines initialization
	timeout := time.After(time.Duration(a.InitTimeout) * time.Second)
	for i:=0; i < len(a.Engines); i++ {
		select{
			// wait to receive an id to set as initialized the engine
			case  <-a.initCh:
				//env.Output.WriteChInfo("(Agent::Start) ready",a.Engines[id].GetName())

			// wait to receive an error
			case err = <-a.initErrCh:
				//env.Output.WriteChError("(Agent::Start)",err)

			// define new error when timeout is reached an not all engines are initialized
			case <- timeout:
				msg := "(Agent::initialize) Not all engines have been initialized after "+strconv.Itoa(int(a.InitTimeout))+" seconds."
				err = errors.New(msg)
			}
	}
	// return when an error is detected after engines' initialization
	if err != nil {
		return err
	}

	return nil
}

// method to run engines
func (a *Agent) run() error {
	var err error
	// runVector contains the engines order to be run
	runVector := []uint{}

	// Set engines on run vector keeping the user defined order
	for _,id := range a.RunOrder {
		runVector = append(runVector, uint(id))
	}
	// Append to run vector engines with no order requiered
	if len(runVector) != len(a.Engines) {
		for engine_id,_ := range a.Engines {
			found := false
			for _,id := range runVector {
				if id == engine_id {
					found = true
				}
			}
			if !found {
				runVector = append(runVector, uint(engine_id))
			}
		}
	}

	// run all engines
	go func(v []uint) {
		var err error
		for _,id := range v {
			e := a.Engines[id]
			err = e.Run()
		}
		if err != nil {
			a.readyErrCh <- err
		} else {
			a.readyCh <- uint(0)
		}
	}(runVector)


	// Waiting for all engines to be run
	// define a timeout to wait for all engines initialization
	timeout := time.After(time.Duration(a.ReadyTimeout) * time.Second)
	select{
		// receive a ready
		case <-a.readyCh:
			// wait to receive an error
		case err = <-a.readyErrCh:
			env.Output.WriteChError("(Agent::run)",err)
		// define new error when timeout is reached an not all engines are initialized
		case <- timeout:
			msg := "(Agent::run) Not all engines have been ready after "+strconv.Itoa(int(a.ReadyTimeout))+" seconds."
			err = errors.New(msg)
			//env.Output.WriteChError(err)
	}

	if err != nil {
		return err
	}

	return nil
}

// setEnginesSubscriptions
func (a *Agent) setEnginesSubscriptions() error {
	for _,e := range a.Engines {
		for _,id := range e.GetDependencies() {
			d := a.Engines[id]
			//if a.EngineStatus[d.GetID()] < engine.INITIALIZED {
			if d.GetStatus() < engine.INITIALIZED {
				return errors.New("(Agent::setEnginesSubscriptions) Engine '"+d.GetName()+"' is not on '"+engine.ToHummanStatus(engine.INITIALIZED)+"' status")
			}
			d.Subscribe(e.GetInputChannel(),e.GetName())
		}
	}
	return nil
}

// validateGraphEngine
func (a *Agent) validateGraphEngine() error {
	var markedEngines map[uint]bool
	for _,e := range a.Engines {
		markedEngines = map[uint]bool{}
		markedEngines[uint(e.GetID())] = true
		for _,id := range e.GetDependencies() {
			if err := a.validateGraphEngineHelper(uint(id), markedEngines); err != nil {
				return err
			}
		}
	}
	return nil
}

// validateGraphEngineHelper
func (a *Agent) validateGraphEngineHelper(id uint, markedEngines map[uint]bool) error {
	var ok bool
	var e engine.Engine

	if e, ok = a.Engines[uint(id)];!ok {
		msg := "(Agent::validateGraphEngineHelper) Unexisten engine with ID '"+fmt.Sprint(id)+"'."
		return errors.New(msg)
	}
	dep := e.GetDependencies()

	if len(dep) == 0 {
		return nil
	} else {
		markedEngines[uint(e.GetID())] = true
		for _,dep_id := range dep {
			if _,ok := markedEngines[uint(dep_id)]; ok {
				msg := "(Agent::validateGraphEngineHelper) Dependency loop with engine '"+e.GetName()+"' and engine '"+a.Engines[uint(dep_id)].GetName()+"'."
				return errors.New(msg)
			} else {
				if err := a.validateGraphEngineHelper(uint(dep_id), markedEngines); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func ToHummanExitStatus(e uint) string {
	exit := map[uint] string {
		OK: "OK",
		WARN: "WARN",
		ERROR: "ERROR",
	}

	humman, ok := exit[e]
	if !ok {
		return "UNKNOWN"
	}
	return humman
}
