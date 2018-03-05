package agent

import (
	"errors"
	"testing"

	"verdmell/engine"

	"github.com/stretchr/testify/assert"
)

func TestValidateGraphEngine(t *testing.T) {
	mock_id_0 := uint(0)
	mock_id_1 := uint(1)
	mock_id_2 := uint(2)

	tests := []struct{
		desc string
		loglevel int
		configdir string
		configfile string
		engines map[uint] engine.Engine
		err error
		InitTimeout int
	}{
		{
			desc: "Testing a basic engine graph.",
			loglevel: 3,
			configdir: "../test/conf.d",
			configfile: "",
			engines: map[uint] engine.Engine {
				mock_id_0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID: mock_id_0,
						Name: "Mock 0",
						Dependencies: []uint{},
					},
				},
				mock_id_1: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID: mock_id_1,
						Name: "Mock 1",
						Dependencies: []uint{mock_id_0},
					},
				},
				mock_id_2: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID: mock_id_1,
						Name: "Mock 2",
						Dependencies: []uint{mock_id_0,mock_id_1},
					},
				},
			},
			err: nil,
		},
		{
			desc: "Testing an engine graph with loops.",
			loglevel: 1,
			configdir: "../test/conf.d",
			configfile: "",
			engines: map[uint] engine.Engine {
				mock_id_0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID: mock_id_0,
						Name: "Mock 0",
						Dependencies: []uint{mock_id_1},
					},
				},
				mock_id_1: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID: mock_id_1,
						Name: "Mock 1",
						Dependencies: []uint{mock_id_0},
					},
				},
			},
			err: errors.New("(Agent::validateGraphEngineHelper) Dependency loop with engine 'Mock 1' and engine 'Mock 0'."),
		},
		{
			desc: "Testing a complex loop dependency.",
			loglevel: 1,
			configdir: "../test/conf.d",
			configfile: "",
			engines: map[uint] engine.Engine {
				mock_id_0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID: mock_id_0,
						Name: "Mock 0",
						Dependencies: []uint{mock_id_1},
					},
				},
				mock_id_1: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID: mock_id_1,
						Name: "Mock 1",
						Dependencies: []uint{mock_id_2},
					},
				},
				mock_id_2: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID: mock_id_1,
						Name: "Mock 2",
						Dependencies: []uint{mock_id_0},
					},
				},
			},
			err: errors.New("(Agent::validateGraphEngineHelper) Dependency loop with engine 'Mock 2' and engine 'Mock 0'."),
		},
		{
			desc: "Testing an unexistent dependency.",
			loglevel: 1,
			configdir: "../test/conf.d",
			configfile: "",
			engines: map[uint] engine.Engine {
				mock_id_0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID: mock_id_0,
						Name: "Mock 0",
						Dependencies: []uint{mock_id_1},
					},
				},
			},
			err: errors.New("(Agent::validateGraphEngineHelper) Unexisten engine with ID '1'."),
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		a := &Agent{
			Loglevel: test.loglevel,
			Configdir: test.configdir,
			Configfile: test.configfile,
			Engines: test.engines,
			InitTimeout: test.InitTimeout,
		}

		err := a.validateGraphEngine()
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		}
	}
}


func TestStart(t *testing.T){
	mock_id_0 := uint(0)
	mock_id_1 := uint(1)

	tests := []struct{
		desc string
		loglevel int
		configdir string
		configfile string
		engines map[uint] engine.Engine
		err error
		InitTimeout int
	}{
		{
			desc: "Testing an agent with no engines.",
			loglevel: 1,
			configdir: "../test/conf.d",
			configfile: "",
			engines: nil,
			err: nil,
		},
		{
			desc: "Testing engine init timeout.",
			loglevel: 1,
			configdir: "../test/conf.d",
			configfile: "",
			InitTimeout: 1,
			engines: map[uint] engine.Engine {
				mock_id_0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID: mock_id_0,
						Name: "Mock 0",
						Dependencies: []uint{},
					},
					InitSleep: 5,
				},
			},
			err: errors.New("(Agent::Start) Not all engines have been initialized after 1 seconds."),
		},
		{
			desc: "Testing an agent with multiple engines.",
			loglevel: 1,
			configdir: "../test/conf.d",
			configfile: "",
			engines: map[uint] engine.Engine {
				mock_id_0: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID: mock_id_0,
						Name: "Mock 0",
						Dependencies: []uint{},
					},
				},
				mock_id_1: &engine.MockEngine{
					BasicEngine: engine.BasicEngine{
						ID: mock_id_1,
						Name: "Mock 1",
						Dependencies: []uint{mock_id_0},
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		a := &Agent{
			Loglevel: test.loglevel,
			Configdir: test.configdir,
			Configfile: test.configfile,
			Engines: test.engines,
			InitTimeout: test.InitTimeout,
		}

		_,err := a.Start()
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		}
		//a.Status()
	}
}