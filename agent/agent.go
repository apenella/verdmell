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
	Engines map[uint]engine.Engine `json: "-"`
	EngineStatus map[uint] uint `json: "-"`

	initCh chan uint
	initErrCh chan error

	InitTimeout int
}

//
// Common methods
//---------------------------------------------------------------------

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

	//initialize engine status datastructure
	a.EngineStatus = make(map[uint] uint)

	// initialize channel for
	a.initCh = make(chan uint)
	a.initErrCh = make(chan error)

	//
	if a.InitTimeout <= 0 {
		a.InitTimeout = 10
	}

	return nil
}

//	Start
func (a *Agent) Start() error {
	err := a.init()
	if err != nil {
		return err
	}

	env.Output.WriteChDebug("(Agent::Start)")

	// Initialize Engines
	for id, e := range a.Engines {
		env.Output.WriteChInfo(id,"=>",e)
		a.EngineStatus[id] = engine.INITIALIZING
		
		go func() {
			err := e.Init()
			if err != nil {
				a.initErrCh <- err
			}
			a.initCh <- id
		}()
	}

	// Waiting for all engines to be initialized
	// 
	timeout := time.After(time.Duration(a.InitTimeout) * time.Second)
	for i:=0; i < len(a.Engines); i++ {
		select{
		case err := <-a.initErrCh: 
			return err
		case id := <-a.initCh: 
			a.EngineStatus[id] = engine.INITIALIZED
		case <- timeout:
			msg := "(Agent::Start) Not all engines have been initialized after "+strconv.Itoa(int(a.InitTimeout))+" seconds."
			return errors.New(msg)
		}
	}
	//if err := a.validateGraphEngine(); err != nil {
	//	return err
	//}

	a.Status()

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


// Status
func (a *Agent) Status() error {
	env.Output.WriteChDebug("(Agent::Status)")

	lines := []string{
		"ENGINE | STATUS",
	}

	for id, e := range a.Engines {
			env.Output.WriteChInfo(e.GetName())
		lines = append(lines, e.GetName()+" | "+engine.ToHummanStatus(a.EngineStatus[id]))
	}

	fmt.Println(columnize.SimpleFormat(lines))

	return nil
}

// String method transform the Configuration to string
func (a *Agent) String() string {
	if err, str := utils.ObjectToJsonString(a); err != nil{
		return err.Error()
	} else{
		return str
	} 
}