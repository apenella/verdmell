package engine

import (
//	"verdmell/environment"
)

/*
	Constants to define engine state
*/
const (
	INITIALIZING = iota
	INITIALIZED
	STARTING
	STARTED
	STOPPING
	STOPPED
	WAITING_DEPENDENCIES
)

/*
	Constants to define engines IDs
*/
const (
	CHECK = iota
	SERVICE
	SAMPLE
	CLUSTER
	API
	UI
	CLIENT
)

/*
	Interface which defines what is an engine
*/
type Engine interface {
	Init() error
	Run() error
	Stop() error
	Subscribe(o chan interface{}, desc string) error
	Status() error
	SayHi()
}