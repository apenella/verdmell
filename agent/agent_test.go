package agent

import (
	"testing"
)

func TestStart(t *testing.T){
	a := &Agent{
		Loglevel: 1,
		Configdir: "../test/conf.d",
		Engines: nil,
	}

	err := a.Start()
	if err != nil {
		t.Fatalf("(Agent::TestStart) ",err)		
	}

}

func TestStatus(t *testing.T) {
	
	a := &Agent{
		Loglevel: 1,
		Configdir: "../test/conf.d",
		Engines: nil,
	}

	err := a.Status()
	if err != nil {
		t.Fatalf("(Agent::TestStatus) ",err)		
	}
}