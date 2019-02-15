/*
Package schedule is used to execute and plan executions
*/
package schedule

// OnceScheduler runs all scheduled task only once
type OnceScheduler struct {
	BasicScheduler
}

// Schedule method append a new unit to the sched
func (s *OnceScheduler) Schedule() {
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

	for _, j := range s.Jobs {
		go j.Unit.Runner.Run()
	}
}

// Stop method append a new unit to the sched
func (s *OnceScheduler) Stop() {
	return
}
