package agent

import (

	"errors"
	"fmt"
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

var env *environment.Environment

// Agent is the element that coordinates all components
type Agent struct{
	// loglevel for agent
	Loglevel int `json: "loglevel"`
	// configuration file name
	Configfile string `json: "configuration_file"`
	// folder to place configuration
	Configdir string `json: "configuration_dir"`

	// List of engines
	Engines map[uint]engine.Engine `json: "engines"`
	EngineStatus map[uint]uint `json: "-"`

	// Timeout to waiting for engines initialization
	InitTimeout int
	// Timeout to waiting for engines to be stopped
	StopTimeout int

	// channels to notify init
	initCh chan uint
	initErrCh chan error

	// channels to notify stop
	stopCh chan uint
	stopErrCh chan error
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

	// Initialize Engines
	for id, e := range a.Engines {
		
		a.EngineStatus[id] = engine.INITIALIZING
		
		// initialize current engine
		go func(id uint, e engine.Engine) {
			err := e.Init()
			// write on error init channel whether there is an error on Init
			if err != nil {
				env.Output.WriteChInfo("(Agent::Start)",err)
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
			case id := <-a.initCh:
				env.Output.WriteChInfo("(Agent::Start) ready",a.Engines[id].GetName())
				a.EngineStatus[id] = engine.READY

			// wait to receive an error
			case err = <-a.initErrCh:
				env.Output.WriteChError("(Agent::Start)",err)
			
			// define new error when timeout is reached an not all engines are initialized
			case <- timeout:
				msg := "(Agent::Start) Not all engines have been initialized after "+strconv.Itoa(int(a.InitTimeout))+" seconds."
				err = errors.New(msg)
				env.Output.WriteChError(err)
			}
	}

	// return when an error is detected after engines' initialization
	if err != nil {
		// stop the agent before return 
		a.Stop()
		return ERROR, err
	}

	if err = a.setEnginesSubscriptions(); err != nil {
		// stop the agent before return 
		a.Stop()
		return ERROR, err		
	}

	// Run engines
	for id,_ := range a.Engines {
		a.EngineStatus[id] = engine.STARTING
	}

	return OK, nil
}

func (a *Agent) Stop() int {
	var err error

	// Initialize Engines
	for id, e := range a.Engines {
		
		a.EngineStatus[id] = engine.STOPPING
		
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
				a.EngineStatus[id] = engine.STOPPED

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
	defer close(a.stopCh)
	defer close(a.stopErrCh)

	// return when an error is detected stopping engines
	if err != nil {
		return ERROR
	}

	a.Status()
	return OK
}

// Status
func (a *Agent) Status() int {
	env.Output.WriteChDebug("(Agent::Status)")

	lines := []string{
		"ENGINE | STATUS",
		"------ | ------",
	}

	for id, e := range a.Engines {
		lines = append(lines, e.GetName()+" | "+engine.ToHummanStatus(a.EngineStatus[id]))
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
//

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

	// initialize engine status datastructure
	a.EngineStatus = make(map[uint] uint)

	// initialize channel to receive just initialited engines id
	a.initCh = make(chan uint)
	// initialize channel to receive error detected during initialization
	a.initErrCh = make(chan error)

	// initialize channel to receive just stopped engines id
	a.stopCh = make(chan uint)
	// initialize channel to receive error when stopping
	a.stopErrCh = make(chan error)

	// set initialization timeout
	if a.InitTimeout <= 0 {
		a.InitTimeout = 30
	}

	// set initialization timeout
	if a.StopTimeout <= 0 {
		a.StopTimeout = 30
	}

	return nil
}

// setEnginesSubscriptions
func (a *Agent) setEnginesSubscriptions() error {
	for _,e := range a.Engines {
		for _,id := range e.GetDependencies() {
			d := a.Engines[id]
			if a.EngineStatus[d.GetID()] != engine.READY {
				return errors.New("(Agent::setEnginesSubscriptions) Engine '"+d.GetName()+"' is not on '"+engine.ToHummanStatus(engine.READY)+"' status")
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