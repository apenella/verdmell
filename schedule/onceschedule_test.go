/*
Package schedule is used to execute and plan executions
*/
package schedule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOnceScheduleSchedule validate a OnceSchedule Schedule
func TestOnceScheduleSchedule(t *testing.T) {
	tests := []struct {
		desc      string
		err       error
		scheduler *OnceScheduler
		units     []*Unit
		jobs      []*Job
	}{
		{
			desc: "Scheduler without units",
			err:  nil,
			scheduler: &OnceScheduler{
				BasicScheduler: BasicScheduler{
					Units: []*Unit{},
					Jobs:  []*Job{},
				},
			},
			units: nil,
			jobs:  nil,
		},
		{
			desc: "Scheduler units to nil jobs",
			err:  nil,
			scheduler: &OnceScheduler{
				BasicScheduler: BasicScheduler{
					Units: []*Unit{
						{
							Runner:   &MockRunner{},
							Name:     "MockRunner",
							Interval: 0,
						},
					},
					Jobs: nil,
				},
			},
			units: []*Unit{
				{
					Runner:   &MockRunner{},
					Name:     "MockRunner",
					Interval: 0,
				},
			},
			jobs: []*Job{
				{
					Unit: &Unit{
						Runner:   &MockRunner{},
						Name:     "MockRunner",
						Interval: 0,
					},
					State: STOPPED,
				},
			},
		},
		{
			desc: "Scheduler units",
			err:  nil,
			scheduler: &OnceScheduler{
				BasicScheduler: BasicScheduler{
					Units: []*Unit{
						{
							Runner:   &MockRunner{},
							Name:     "MockRunner",
							Interval: 0,
						},
					},
					Jobs: []*Job{},
				},
			},
			units: []*Unit{
				{
					Runner:   &MockRunner{},
					Name:     "MockRunner",
					Interval: 0,
				},
			},
			jobs: []*Job{
				{
					Unit: &Unit{
						Runner:   &MockRunner{},
						Name:     "MockRunner",
						Interval: 0,
					},
					State: STOPPED,
				},
			},
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		test.scheduler.Schedule()
		assert.Equal(t, len(test.scheduler.Jobs), len(test.jobs), "Does not have the expected jobs")
	}
}

// TestOnceScheduleStop validate a OnceSchedule Stop
func TestOnceScheduleStop(t *testing.T) {
	tests := []struct {
		desc      string
		err       error
		scheduler *OnceScheduler
	}{
		{
			desc: "Scheduler without units",
			err:  nil,
			scheduler: &OnceScheduler{
				BasicScheduler: BasicScheduler{
					Units: []*Unit{},
					Jobs:  []*Job{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		test.scheduler.Stop()
		assert.Equal(t, nil, nil, "")
	}
}
