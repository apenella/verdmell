/*
Package agent
*/
package agent

import (
	"errors"
	"testing"

	"verdmell/engine"

	"github.com/stretchr/testify/assert"
)

func TestValidateGraphEngine(t *testing.T) {
	mockID0 := uint(0)
	mockID1 := uint(1)
	mockID2 := uint(2)

	tests := []struct {
		desc        string
		loglevel    int
		configdir   string
		configfile  string
		engines     map[uint]engine.Engine
		err         error
		InitTimeout int
	}{
		{
			desc:       "Testing a basic engine graph.",
			loglevel:   3,
			configdir:  "../test/conf.d",
			configfile: "",
			engines: map[uint]engine.Engine{
				mockID0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:           mockID0,
						Name:         "Mock 0",
						Dependencies: []uint{},
					},
				},
				mockID1: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:           mockID1,
						Name:         "Mock 1",
						Dependencies: []uint{mockID0},
					},
				},
				mockID2: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:           mockID2,
						Name:         "Mock 2",
						Dependencies: []uint{mockID0, mockID1},
					},
				},
			},
			err: nil,
		},
		{
			desc:       "Testing an engine graph with loops.",
			loglevel:   1,
			configdir:  "../test/conf.d",
			configfile: "",
			engines: map[uint]engine.Engine{
				mockID0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:           mockID0,
						Name:         "Mock 0",
						Dependencies: []uint{mockID1},
					},
				},
				mockID1: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:           mockID1,
						Name:         "Mock 1",
						Dependencies: []uint{mockID0},
					},
				},
			},
			err: errors.New("(Agent::validateGraphEngineHelper) Dependency loop with engine 'Mock 1' and engine 'Mock 0'."),
		},
		{
			desc:       "Testing an unexistent dependency.",
			loglevel:   1,
			configdir:  "../test/conf.d",
			configfile: "",
			engines: map[uint]engine.Engine{
				mockID0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:           mockID0,
						Name:         "Mock 0",
						Dependencies: []uint{mockID1},
					},
				},
			},
			err: errors.New("(Agent::validateGraphEngineHelper) Unexisten engine with ID '1'."),
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		a := &Agent{
			Engines:     test.engines,
			InitTimeout: test.InitTimeout,
		}
		/*
		   TO BE REVIEW

		   agent_test.go:101: Testing an engine graph with loops.
		   	 agent_test.go:109:
		   				 Error Trace:	agent_test.go:109
		   				 Error:      	Not equal:
		   											 expected: &errors.errorString{s:"(Agent::validateGraphEngineHelper) Dependency loop with engine 'Mock 1' and engine 'Mock 0'."}
		   											 actual  : &errors.errorString{s:"(Agent::validateGraphEngineHelper) Dependency loop with engine 'Mock 0' and engine 'Mock 1'."}

		   											 Diff:
		   											 --- Expected
		   											 +++ Actual
		   											 @@ -1,2 +1,2 @@
		   											 -(*errors.errorString)((Agent::validateGraphEngineHelper) Dependency loop with engine 'Mock 1' and engine 'Mock 0'.)
		   											 +(*errors.errorString)((Agent::validateGraphEngineHelper) Dependency loop with engine 'Mock 0' and engine 'Mock 1'.)

		   				 Test:       	TestValidateGraphEngine
		*/
		err := a.validateGraphEngine()
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		}
	}
}

func TestSetEngineSubscriptions(t *testing.T) {
	mockID0 := uint(0)
	mockID1 := uint(1)
	mockID2 := uint(2)

	tests := []struct {
		desc         string
		loglevel     int
		configdir    string
		configfile   string
		engines      map[uint]engine.Engine
		engineStatus map[uint]uint
		err          error
	}{
		{
			desc:       "Testing subscriptions.",
			loglevel:   1,
			configdir:  "../test/conf.d",
			configfile: "",
			engines: map[uint]engine.Engine{
				mockID0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID0,
						Name:          "Mock 0",
						Dependencies:  []uint{},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
						Status:        engine.INITIALIZED,
					},
				},
				mockID1: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID1,
						Name:          "Mock 1",
						Dependencies:  []uint{mockID0},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
						Status:        engine.INITIALIZED,
					},
				},
				mockID2: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID2,
						Name:          "Mock 2",
						Dependencies:  []uint{mockID0, mockID1},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
						Status:        engine.INITIALIZED,
					},
				},
			},
			err: nil,
		},
		{
			desc:       "Testing subscribe to non initialized engine.",
			loglevel:   1,
			configdir:  "../test/conf.d",
			configfile: "",
			engines: map[uint]engine.Engine{
				mockID0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID0,
						Name:          "Mock 0",
						Dependencies:  []uint{},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
						Status:        engine.STOPPED,
					},
				},
				mockID1: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID1,
						Name:          "Mock 1",
						Dependencies:  []uint{mockID0},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
						Status:        engine.READY,
					},
				},
				mockID2: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID2,
						Name:          "Mock 2",
						Dependencies:  []uint{mockID0, mockID1},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
						Status:        engine.READY,
					},
				},
			},
			err: errors.New("(Agent::setEnginesSubscriptions) Engine 'Mock 0' is not on 'Initialized' status"),
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		a := &Agent{
			Engines:      test.engines,
			EngineStatus: test.engineStatus,
		}

		err := a.setEnginesSubscriptions()
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		}
	}
}

// Test run
func TestInitialize(t *testing.T) {
	mockID0 := uint(0)

	tests := []struct {
		desc         string
		loglevel     int
		configdir    string
		configfile   string
		engines      map[uint]engine.Engine
		engineStatus map[uint]uint
		runOrder     []uint
		InitTimeout  int
		err          error
	}{
		{
			desc:        "Testing engine initialize.",
			loglevel:    1,
			configdir:   "../test/conf.d",
			configfile:  "",
			InitTimeout: 30,
			engines: map[uint]engine.Engine{
				mockID0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:           mockID0,
						Name:         "Mock 0",
						Dependencies: []uint{},
					},
				},
			},
			err: nil,
		},
		{
			desc:        "Testing engine initialize timeout.",
			loglevel:    1,
			configdir:   "../test/conf.d",
			configfile:  "",
			InitTimeout: 1,
			engines: map[uint]engine.Engine{
				mockID0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:           mockID0,
						Name:         "Mock 0",
						Dependencies: []uint{},
					},
					InitSleep: 5,
				},
			},
			err: errors.New("(Agent::initialize) Not all engines have been initialized after 1 seconds."),
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		a := &Agent{
			Engines:      test.engines,
			EngineStatus: test.engineStatus,
			RunOrder:     test.runOrder,
			InitTimeout:  test.InitTimeout,
			initCh:       make(chan uint),
			initErrCh:    make(chan error),
		}

		err := a.initializeEngines()
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		}
		//a.Status()
	}

}

// Test run
func TestRun(t *testing.T) {
	mockID0 := uint(0)
	mockID1 := uint(1)
	mockID2 := uint(2)
	mock_id_3 := uint(3)

	tests := []struct {
		desc         string
		loglevel     int
		configdir    string
		configfile   string
		engines      map[uint]engine.Engine
		engineStatus map[uint]uint
		runOrder     []uint
		readyTimeout int
		err          error
	}{
		{
			desc:         "Testing run timout error.",
			loglevel:     3,
			configdir:    "../test/conf.d",
			configfile:   "",
			readyTimeout: 1,
			engines: map[uint]engine.Engine{
				mockID0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID0,
						Name:          "Mock 0",
						Dependencies:  []uint{},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
					},
					RunSleep: 5,
				},
			},
			err: errors.New("(Agent::run) Not all engines have been ready after 1 seconds."),
		},
		{
			desc:         "Testing basic run without run order.",
			loglevel:     1,
			configdir:    "../test/conf.d",
			configfile:   "",
			readyTimeout: 30,
			engines: map[uint]engine.Engine{
				mockID0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID0,
						Name:          "Mock 0",
						Dependencies:  []uint{},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
					},
				},
				mockID1: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID1,
						Name:          "Mock 1",
						Dependencies:  []uint{},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
					},
				},
				mockID2: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID2,
						Name:          "Mock 2",
						Dependencies:  []uint{},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
					},
				},
				mock_id_3: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mock_id_3,
						Name:          "Mock 3",
						Dependencies:  []uint{},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
					},
				},
			},
			err: nil,
		},
		{
			desc:         "Testing basic run with run order.",
			loglevel:     1,
			configdir:    "../test/conf.d",
			configfile:   "",
			readyTimeout: 30,
			engines: map[uint]engine.Engine{
				mockID0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID0,
						Name:          "Mock 0",
						Dependencies:  []uint{},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
					},
				},
				mockID1: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID1,
						Name:          "Mock 1",
						Dependencies:  []uint{},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
					},
				},
				mockID2: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mockID2,
						Name:          "Mock 2",
						Dependencies:  []uint{},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
					},
				},
				mock_id_3: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:            mock_id_3,
						Name:          "Mock 3",
						Dependencies:  []uint{},
						Subscriptions: make(map[chan interface{}]string),
						InputChannel:  make(chan interface{}),
					},
				},
			},
			runOrder: []uint{
				mockID2,
				mock_id_3,
				mockID0,
				mockID1,
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		a := &Agent{
			Engines:      test.engines,
			EngineStatus: test.engineStatus,
			RunOrder:     test.runOrder,
			ReadyTimeout: test.readyTimeout,
			readyCh:      make(chan uint),
			readyErrCh:   make(chan error),
		}

		err := a.runEngines()
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		}
		//a.Status()
	}
}

// Start
func TestStart(t *testing.T) {
	mockID0 := uint(0)
	mockID1 := uint(1)

	tests := []struct {
		desc         string
		loglevel     int
		configdir    string
		configfile   string
		engines      map[uint]engine.Engine
		err          error
		InitTimeout  int
		ReadyTimeout int
	}{
		{
			desc:         "Testing an agent with no engines.",
			loglevel:     1,
			configdir:    "../test/conf.d",
			configfile:   "",
			engines:      nil,
			err:          nil,
			ReadyTimeout: 60,
		},
		{
			desc:       "Testing an agent with multiple engines.",
			loglevel:   1,
			configdir:  "../test/conf.d",
			configfile: "",
			engines: map[uint]engine.Engine{
				mockID0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:           mockID0,
						Name:         "Mock 0",
						Dependencies: []uint{},
					},
				},
				mockID1: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID:           mockID1,
						Name:         "Mock 1",
						Dependencies: []uint{mockID0},
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		a := &Agent{
			Context:      nil,
			Engines:      test.engines,
			InitTimeout:  test.InitTimeout,
			ReadyTimeout: test.ReadyTimeout,
		}

		_, err := a.Start()
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		}
	}
}
