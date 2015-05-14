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
  //Get defined checks
  ck := RetrieveChecks(folder)

  if err = ck.ValidateChecks(nil); err == nil {
    cks.SetChecks(ck)
    //Init the running queues to proceed the executions
    cks.InitCheckRunningQueues()
  } else {
    return err, nil
  }

   if err, ss = sample.NewSampleSystem(env); err == nil {
     cks.SetSampleSystem(ss)
   } else {
     return err, nil
   }

  //Get defined checks groups
  cg := RetrieveCheckgroups(folder)
  if err := cg.ValidateCheckgroups(ck); err == nil {
    cks.SetCheckgroups(cg)
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

  switch req := i.(type){
  case *CheckObject:
    env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Starting the check: "+req.String())
    
    // starting the check object task for the gived check
    //_,result = req.StartCheckObjectTask()
    
    // run a goroutine for each checkObject and write the result to the channel
    go func() {
      _,res := req.StartCheckObjectTask();
      statusChan <- res
    }()
    // waiting the CheckObject result
    exitStatus = <-statusChan
    env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Check '"+strconv.Itoa(exitStatus)+"' done")

  case []string:
    env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Running a Checkgroup")

    for _,checkname := range req {
      env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Running the checkgroup check: "+checkname)
      checks := c.GetChecks()

      if checkObj, err := checks.GetCheckObjectByName(checkname); err != nil {
        return err,2
      } else {        

        // run a goroutine for each checkObject and write the result to the channel
        go func() {
          _,res := checkObj.StartCheckObjectTask();
          statusChan <- res
        }()
      }
    }
    // waiting the CheckObjects results
    for i:= 0; i<len(req); i++{
      subExitStatus := <-statusChan
      if exitStatus < subExitStatus {
        exitStatus = subExitStatus
      }
      env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Check '"+strconv.Itoa(subExitStatus)+"' done")
    }
    
  default:
    env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Running all checks")
    
    // Get Checks attribute from CheckSystem
    checks :=  c.GetChecks()
    _,exitStatus = checks.StartCheckTaskPools()
  }

  if exitStatus < 0 {
    exitStatus = 3
  }

  return nil,exitStatus
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
