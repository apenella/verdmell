package engine

import (
	"sync"

	"verdmell/context"
)

// Engine is an interface to be accomplish for any verdmell component.

// Constants to define engine status
const (
	INITIALIZING = iota
	INITIALIZED
	STARTING
	READY
	STOPPING
	STOPPED
	WAITING_DEPENDENCIES
	NOT_INITIALIZED
)

// Constants to define engines IDs for known engines
const (
	CHECK = iota
	SERVICE
	SAMPLE
	CLUSTER
	API
	UI
	CLIENT
)

// Engine interface which defines what is an engine
type Engine interface {
	GetID() uint
	GetName() string
	GetDependencies() []uint
	GetInputChannel() chan interface{}
	GetStatus() uint
	SetStatus(s uint)

	Init() error
	Run() error
	Stop() error
	Subscribe(o chan interface{}, desc string) error
	SayHi()
}

// ToHummanStatus translates status to humman readable one
func ToHummanStatus(s uint) string {

	status := map[uint]string{
		INITIALIZING:         "Initializing",
		INITIALIZED:          "Initialized",
		STARTING:             "Starting",
		READY:                "Ready",
		STOPPING:             "Stopping",
		STOPPED:              "Stopped",
		WAITING_DEPENDENCIES: "Waiting dependencies",
		NOT_INITIALIZED:      "Not initialized",
	}

	humman, ok := status[s]
	if !ok {
		return "Unknown"
	}
	return humman
}

//
// Basic engine definition
type BasicEngine struct {
	ID           uint   `json: "id"`
	Name         string `json: "name"`
	Dependencies []uint `json: "dependencies"`
	// Context contains information about the runtime state
	Context *context.Context
	// subscriptions Channel
	Subscriptions map[chan interface{}]string `json: "-"`
	// input channel
	InputChannel chan interface{} `json: "-"`
	// engine's current status
	Status uint `json: "status"`
	// status mutex
	statusMutex sync.RWMutex `json: "-"`
	// init mutex
	initMutex sync.RWMutex `json: "-"`
	// run mutex
	runMutex sync.RWMutex `json: "-"`
	// stop mutex
	stopMutex sync.RWMutex `json: "-"`
}

// GetID
func (e *BasicEngine) GetID() uint {
	return e.ID
}

// GetName
func (e *BasicEngine) GetName() string {
	return e.Name
}

// GetDependencies
func (e *BasicEngine) GetDependencies() []uint {
	return e.Dependencies
}

// GetStatus
func (e *BasicEngine) GetStatus() uint {
	//read lock
	e.statusMutex.RLock()
	defer e.statusMutex.RUnlock()

	return e.Status
}

// SetStatus
func (e *BasicEngine) SetStatus(s uint) {
	//write lock
	e.statusMutex.Lock()
	defer e.statusMutex.Unlock()

	e.Status = s
}

// SetSubscriptions
func (e *BasicEngine) SetSubscriptions(s map[chan interface{}]string) {
	e.Subscriptions = s
}

// GetSubscriptions
func (e *BasicEngine) GetSubscriptions() map[chan interface{}]string {
	return e.Subscriptions
}

// SetSubscriptions
func (e *BasicEngine) SetInputChannel(c chan interface{}) {
	e.InputChannel = c
}

// GetInputChannel
func (e *BasicEngine) GetInputChannel() chan interface{} {
	return e.InputChannel
}
