package command


//	ExecCommand
type HealthCheckCommand struct{}

//	Run
func (c *HealthCheckCommand) Run(args []string) int {
  return 0
}

// Help
func (c *HealthCheckCommand) Help() string {
	return "Usage: verdmell exec [options]"
}

// Synopsis
func (c *HealthCheckCommand) Synopsis() string {
	return "Execute checks on isolated mode"
}
