/*
Check system management
-Checks
-Check groups
-CheckSamplesMap
-CheckSamples
-RunningQueues
-CheckQueue
*/
package check

import (
  "strconv"
  "verdmell/environment"
  "verdmell/sample"
)
//
var env *environment.Environment
//#
//#
//# CheckSystem struct
//# The struct for CheckSystem has all Check, Checkgroup and samples information
type CheckSystem struct{
  // Map to storage the checks
  Ck *Checks
  // Map to storage the checkgroups
  Cg *Checkgroups
  // Map to storage the samples
  Cs *sample.SampleSystem
}
//
//# NewCheckSystem: return a Checksystem instance to be run
func NewCheckSystem(e *environment.Environment) (error, *CheckSystem){
  e.Output.WriteChDebug("(CheckSystem::NewCheckSystem)")
  cks := new(CheckSystem)
  ss := new(sample.SampleSystem)
  var err error

  // get the environment attributes
  env = e

  // folder contains check definitions
  folder := env.Setup.Checksfolder
  // Get defined checks
  // validate checks and set the checks into check system
  ck := RetrieveChecks(folder)
  if err = ck.ValidateChecks(nil); err == nil {
    cks.SetChecks(ck)
    //Init the running queues to proceed the executions
    cks.InitCheckRunningQueues()
  } else {
    return err, nil
  }

  // Get defined checks groups
  // validate checks and set the checks into check system
  cg := RetrieveCheckgroups(folder)
  if err := cg.ValidateCheckgroups(ck); err == nil {
    cks.SetCheckgroups(cg)
  } else {
    return err, nil
  }  
  
  // starting sample system
  if err, ss = sample.NewSampleSystem(env); err == nil {
   cks.SetSampleSystem(ss)
  } else {
   return err, nil
  }

	return err, cks
}
//#
//# Getters and Setters
//#----------------------------------------------------------------------------------------

//
//# SetChecks: attribute from CheckSystem
func (c *CheckSystem) SetChecks(ck *Checks) {
  env.Output.WriteChDebug("(CheckSystem::SetChecks) Set value '"+ck.String()+"'")
  c.Ck = ck
}
//
//# SetCheckgroups: attribute from CheckSystem
func (c *CheckSystem) SetCheckgroups(cg *Checkgroups) {
  env.Output.WriteChDebug("(CheckSystem::SetCheckgroups) Set value '"+cg.String()+"'")
  c.Cg = cg
}
//
//# SetChecksamplesmap: attribute from CheckSystem
func (c *CheckSystem) SetSampleSystem(cs *sample.SampleSystem) {
  env.Output.WriteChDebug("(CheckSystem::SetSampleSystem)")
  c.Cs = cs
}
//
//# Getchecks: attribute from CheckSystem
func (c *CheckSystem) GetChecks() *Checks{
  env.Output.WriteChDebug("(CheckSystem::GetChecks) Get value")
  return c.Ck
}
//
//# Getcheckgroups: attribute from CheckSystem
func (c *CheckSystem) GetCheckgroups() *Checkgroups{
  env.Output.WriteChDebug("(CheckSystem::GetCheckgroups) Get value")
  return c.Cg
}
//
//# GetChecksamplesmap: attribute from CheckSystem
func (c *CheckSystem) GetSampleSystem() *sample.SampleSystem{
  env.Output.WriteChDebug("(CheckSystem::SetSampleSystem)")
  return c.Cs
}
//
//# StartCheckSystem: will determine which kind of check has been required by user and start the checks
func (c *CheckSystem) StartCheckSystem(i interface{}) (error,int) {
  env.Output.WriteChDebug("(CheckSystem::StartCheckSystem)")
  exitStatus := -1
  statusChan := make(chan int)
  defer close(statusChan)

  // the next will be only used during ec or eg
  // check will contain the Check configurations
  check := new(Checks)
  // checks will contain all the CheckObject definition
  checks := make(map[string]CheckObject)
  //

  switch req := i.(type){
  case *CheckObject:
    env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Starting the check '"+req.String()+"'")
    //add the check to be executed
    checks[req.GetName()] = *req
    //add the check dependencies
    for _,dependency := range req.GetDepend(){
      if checkObj, err := c.Ck.GetCheckObjectByName(dependency); err != nil {
        return err,2
      } else {
        if _,exist := checks[dependency]; !exist{
          checks[dependency] = *checkObj
        }
      }
    }

    check.SetCheck(checks)

    // run a goroutine for each checkObject and write the result to the channel
    go func() {
      _,res := check.StartCheckTaskPools(c.GetSampleSystem())
      statusChan <- res
    }()
    // waiting the CheckObject result
    exitStatus = <-statusChan
    env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Check '"+strconv.Itoa(exitStatus)+"' done")

  case []string:
    env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Running a Checkgroup")

    for _,checkname := range req {
      env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Preparing the check '"+checkname+"'")
      cks := c.GetChecks()

      if checkObj, err := cks.GetCheckObjectByName(checkname); err != nil {
        return err,2
      } else {
        if _,exist := checks[checkname]; !exist{
          checks[checkname] = *checkObj
        }
        //add the check dependencies
        for _,dependency := range checkObj.GetDepend(){
          if checkObjdependency, err := c.Ck.GetCheckObjectByName(dependency); err != nil {
            return err,2
          } else {
            if _,exist := checks[dependency]; !exist{
              checks[dependency] = *checkObjdependency
            }
          }
        }
      }
    }
    check.SetCheck(checks)
    // run a goroutine for each checkObject and write the result to the channel
    go func() {
      _,res := check.StartCheckTaskPools(c.GetSampleSystem())
      statusChan <- res
    }()

    // waiting the CheckObjects results
    exitStatus = <-statusChan
    env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Check '"+strconv.Itoa(exitStatus)+"' done")
    // for i:= 0; i<len(req); i++{
    //   subExitStatus := <-statusChan
    //   if exitStatus < subExitStatus {
    //     exitStatus = subExitStatus
    //   }
    //   env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Check '"+strconv.Itoa(subExitStatus)+"' done")
    // }
    
  default:
    //env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Running all checks")
    // Get Checks attribute from CheckSystem
    //checks :=  c.GetChecks()
    //_,exitStatus = checks.StartCheckTaskPools()
    _,exitStatus = c.GetChecksExitStatus()
  }

  if exitStatus < 0 {
    exitStatus = 3
  }

  return nil,exitStatus
}

//
//# GetChecksExitStatus: return the status of all checks
func (c *CheckSystem) GetChecksExitStatus() (error, int) {
  env.Output.WriteChDebug("(CheckSystem::GetChecksExitStatus) Running all checks")
  // Get Checks attribute from CheckSystem
  checks :=  c.GetChecks()

  env.Output.WriteChDebug("(CheckSystem::GetChecksExitStatus) Samples:"+c.Cs.String())

  _,exitStatus := checks.StartCheckTaskPools(c.GetSampleSystem())

  c.GetChecksAllSamples()

  return nil, exitStatus
}

//
//# GetChecksSamples: return the status of all checks
func (c *CheckSystem) GetChecksAllSamples() {
  env.Output.WriteChDebug("(CheckSystem::GetChecksSamples)")
  // Get Checks attribute from CheckSystem
  checks :=  c.GetChecks()
  samplesystem := c.GetSampleSystem()

  for check := range checks.Check {
    _,s := samplesystem.GetSample(check)
    env.Output.WriteChDebug("(CheckSystem::GetChecksSamples) "+s.String())
  }
}

//
//# InitCheckRunningQueues: prepares each checkobject to be run
func (c *CheckSystem) InitCheckRunningQueues() error {
  cs := c.GetChecks()

  for _,obj := range cs.GetCheck() {
      go func(checkObj CheckObject) {
        env.Output.WriteChDebug("(CheckSystem::InitCheckRunningQueues) CheckQueue for '"+checkObj.GetName()+"'") 
        checkObj.StartQueue()
      }(obj)
  }
  return nil
}
// InitCheckRunningQueues
func (c *CheckSystem) TestCheckRunningQueues() error {
  cs := c.GetChecks()

  for _,obj := range cs.GetCheck() {
      env.Output.WriteChDebug("(CheckSystem::TestCheckRunningQueues) CheckQueue for '"+obj.GetName()+"'")
      env.Output.WriteChDebug(obj)
      obj.EnqueueCheckObject()
  }  
  return nil
}
//#######################################################################################################