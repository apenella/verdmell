package engine

import (
	"errors"
	"fmt"
	"time"
)

// MockEngine is an implementation of Engine that can be used for tests.
type MockEngine struct {
	BasicEngine

	// if its value is greater than 0, the init method gets slept value seconds
	InitSleep int
	// if its value is greater than 0, the run method gets slept value seconds
	RunSleep int
}

//
// functions required to implement an engine
//

// Init
func (m *MockEngine)Init() error {

	m.SetStatus(INITIALIZING)
	// sleep when InitSleep value is greater than 0
	if m.InitSleep > 0 {
		time.Sleep(time.Duration(m.InitSleep) * time.Second)
	}
	
	// define engine name unless engine has already a name
	if m.Name == "" {
		m.Name = "UnknownEngine"
	}

	// initialize subscriptions
	m.SetSubscriptions(make(map[chan interface{}] string))
	// initialize channel to receive notifications
	m.SetInputChannel(make(chan interface{}))

	m.SetStatus(INITIALIZED)
	return nil
}

// Run
func (m *MockEngine)Run() error {
	m.SetStatus(STARTING)
	fmt.Println("(MockEngine::Run) ("+fmt.Sprint(m.ID)+") "+m.Name)
	
	// sleep when RunSleep value is greater than 0
	if m.RunSleep > 0 {
		time.Sleep(time.Duration(m.RunSleep) * time.Second)
	}

	m.SetStatus(READY)
	return nil
}

// Stop
func (m *MockEngine)Stop() error {
	fmt.Println("(MockEngine::Stop) ("+fmt.Sprint(m.ID)+") "+m.Name)
	
	if m.InputChannel != nil {
		defer close(m.InputChannel)
	}

	return nil
}

// Subscribe
func (m *MockEngine)Subscribe(o chan interface{}, desc string) error {
	fmt.Println("(MockEngine::Subscribe) ("+fmt.Sprint(m.ID)+") "+m.Name)

	channels := m.GetSubscriptions()
	if _, exist := channels[o]; !exist {
		channels[o] = desc
	} else {
		return errors.New("(MockEngine::Subscribe) You are trying to add an existing channel")
	}

	return nil
}

// SayHi
func (m *MockEngine)SayHi() {
	fmt.Println("(MockEngine::SayHi) Hi! I'm a MockEngine")
}

//
// Common methods
//

func (m *MockEngine)String() string {
	str := "{"
	str += " ID:"
	str += fmt.Sprint(m.ID)
	str += " Name:"
	str += m.Name
	str += " Dependencies:"
	for _,d := range m.GetDependencies(){
		str += fmt.Sprint(d)
	}
	str += "}"
	return str
	//return "{ ID: "+fmt.Sprint(m.ID)+", Name: "+m.Name+" }"
}