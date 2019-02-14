/*
Package schedule is used to execute and plan executions
*/
package schedule

// OnceScheduler runs all scheduled task only once
type OnceScheduler struct {
	BasicScheduler
}

// Schedule start the tasks to be executed
func (s *OnceScheduler) Schedule() error {
	return nil
}

// Stop
func (s *OnceScheduler) Stop() error {
	return nil
}
