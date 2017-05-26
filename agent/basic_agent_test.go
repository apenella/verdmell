package agent

import (
	"testing"
)

func TestStart(t *testing.T){
	a := &BasicAgent{
		Agent: Agent{
			Loglevel: 1,
			Configdir: "../test/conf.d",
		},
	}

	err := a.Start()
	if err != nil {
		t.Fatalf("(BasicAgent::TestStart)",err)		
	}

}

func TestStatus(t *testing.T) {
	
	a := &BasicAgent{
		Agent: Agent{
			Loglevel: 1,
			Configdir: "../test/conf.d",
		},
	}

	err := a.Status()
	if err.Error() != "(BasicAgent::Status) Null CheckEngine" {
		t.Fatalf("(BasicAgent::TestStatus)",err)
	}

}