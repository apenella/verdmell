/*
Package schedule is used to execute and plan executions
*/
package schedule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAdd validate a BasicSchedule
func TestAdd(t *testing.T) {
	tests := []struct {
		desc      string
		err       error
		scheduler *BasicScheduler
		units     []*Unit
		Unit      *Unit
		size      int
	}{
		{
			desc: "Add a new Unit to new BasicScheduler",
			err:  nil,
			scheduler: &BasicScheduler{
				Units: []*Unit{},
			},
			Unit: &Unit{
				Runner:   &MockRunner{},
				Name:     "MockRunner",
				Interval: 0,
			},
			size: 1,
		},
		{
			desc: "Add a new Unit to nil Units",
			err:  nil,
			scheduler: &BasicScheduler{
				Units: nil,
			},
			Unit: &Unit{
				Runner:   &MockRunner{},
				Name:     "MockRunner",
				Interval: 0,
			},
			size: 1,
		},
		{
			desc: "Add a new Unit to and existing Units BasicScheduler",
			err:  nil,
			scheduler: &BasicScheduler{
				Units: []*Unit{
					&Unit{
						Runner:   &MockRunner{},
						Name:     "First MockRunner",
						Interval: 0,
					},
				},
			},
			Unit: &Unit{
				Runner:   &MockRunner{},
				Name:     "Second MockRunner",
				Interval: 0,
			},
			size: 2,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		test.scheduler.Add(test.Unit)
		assert.Equal(t, len(test.scheduler.Units), test.size, "Does not have the expected units")
	}

}
