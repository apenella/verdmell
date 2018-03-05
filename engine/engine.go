package engine

// Engine is an interface to be accomplish for any verdmell component.

// Constants to define engine status
const (
	INITIALIZING = iota
	READY
	STARTING
	RUNNING
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

// Interface which defines what is an engine
type Engine interface {
	GetID() uint
	GetName() string
	GetDependencies() []uint
	GetInputChannel() chan interface{}
	Init() error
	Run() error
	Stop() error
	Subscribe(o chan interface{}, desc string) error
	Status() int
	SayHi()
}

// ToHummanStatus translates status to humman readable one
func ToHummanStatus(s uint) (string) {

	status := map[uint] string {
		INITIALIZING: "INITIALIZING",
		READY: "READY",
		STARTING: "STARTING",
		RUNNING: "RUNNING",
		STOPPING: "STOPPING",
		STOPPED: "STOPPED",
		WAITING_DEPENDENCIES: "WAITING_DEPENDENCIES",
		NOT_INITIALIZED: "NOT_INITIALIZED",
	}

	humman, ok := status[s]
	if !ok {
		return "UNKNOWN"
	}		
	return humman
}

//
// Basic engine definition
type BasicEngine struct {
	ID uint `json: "id"`
	Name string `json: "name"`
	Dependencies []uint `json: "dependencies"`
	// subscriptions Channel
	subscriptions map[chan interface{}] string `json: "-"`
	// input channel
	inputChannel chan interface{}`json: "-"`
}
// GetID
func (e *BasicEngine)GetID() uint {
	return e.ID
}
// GetName
func (e *BasicEngine)GetName() string {
	return e.Name
}
// GetDependencies
func (e *BasicEngine)GetDependencies() []uint {
	return e.Dependencies
}
// SetSubscriptions
func (e *BasicEngine) SetSubscriptions(s map[chan interface{}] string ) {
	e.subscriptions = s
}
// GetSubscriptions
func (e *BasicEngine) GetSubscriptions() map[chan interface{}] string {
  return e.subscriptions
}
// SetSubscriptions
func (e *BasicEngine) SetInputChannel(c chan interface{}) {
	e.inputChannel = c
}
// GetInputChannel
func (e *BasicEngine) GetInputChannel() chan interface{} {
  return e.inputChannel
}