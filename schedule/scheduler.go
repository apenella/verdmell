/*
Package schedule is used to execute and plan executions
*/
package schedule

import (
	"verdmell/context"
)

// job states
const (
	STOPPED = iota
	RUNNING
	SCHEDULED
)

// Runner interface define those elements which can execute jobs
type Runner interface {
	Run()
}

// Unit is the element which contains all definition for an execution
type Unit struct {
	// Runner is the which could execute the unit
	Runner Runner `json:"runner"`
	// Name is an id of the unit
	Name string `json:"name"`
	// Interval of execution
	Interval int `json:"interval"`
	// Notify is a func which could be used as a callback when a unit finishes to be run
	Notify func() `json:"-"`
}

// Job is and structure which defines a execution Unit in runtime
type Job struct {
	// Unit is the scheduled element to be run
	Unit *Unit `json:"unit"`
	// State of job
	State uint `json:"state"`
}

// Scheduler interface defines and element which could be used to run jobs
type Scheduler interface {
	Schedule() error
	Stop() error
	Add(u *Unit)
}

// SchedulerFactory is a type of function that is a factory for schedulers.
type SchedulerFactory func() (Scheduler, error)

// BasicScheduler is the basic structure of an schedule
type BasicScheduler struct {
	// Context contains information about the whole runtime state
	Context *context.Context `json:"-"`
	// Units has the definition of all execution units
	Units []*Unit `json:"units"`
	// Jobs has the
	Jobs []*Job `json:"jobs"`
}

// Schedule method append a new unit to the sched
func (s *BasicScheduler) Schedule(u *Unit) {
	if s.Jobs == nil {
		s.Jobs = []*Job{}
	}
	for _, u := range s.Units {
		job := &Job{
			Unit:  u,
			State: STOPPED,
		}
		s.Jobs = append(s.Jobs, job)
	}
}

// Stop method append a new unit to the sched
func (s *BasicScheduler) Stop(u *Unit) {
	return
}

// Add method append a new unit to the sched
func (s *BasicScheduler) Add(u *Unit) {
	if s.Units == nil {
		s.Units = []*Unit{}
	}

	s.Units = append(s.Units, u)
}
