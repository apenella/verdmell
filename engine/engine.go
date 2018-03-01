package engine

// Engine is an interface to be accomplish for any verdmell component.

// Constants to define engine status
const (
	INITIALIZING = iota
	INITIALIZED
	STARTING
	STARTED
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
		INITIALIZED: "INITIALIZED",
		STARTING: "STARTING",
		STARTED: "STARTED",
		STOPPING: "STOPPING",
		STOPPED: "STOPPED",
		WAITING_DEPENDENCIES: "WAITING_DEPENDENCIES",
		NOT_INITIALIZED: "NOT_INITIALIZED",
	}

	return status[s]
}

//
// Basic engine definition
type BasicEngine struct {
	ID uint
	Name string
	Dependencies []uint
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