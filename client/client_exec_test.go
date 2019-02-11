package client

import (
	"testing"

	"verdmell/engine"
)

func TestRun(t *testing.T) {

	c := &Client{
		ID: engine.CLIENT,
		Worker: &Exec{
			Checks: []string{"foo", "bar"},
		},
	}

	err := c.Run()
	if err != nil {
		t.Fatalf("(Client::TestRun) ", err)
	}

}
