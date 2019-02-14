/*
Package schedule is used to execute and plan executions
*/
package schedule

// MockRunner implements a Runner
type MockRunner struct{}

// Run does nothing
func (r *MockRunner) Run() {}
