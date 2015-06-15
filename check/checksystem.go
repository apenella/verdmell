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
  "verdmell/environment"
  "verdmell/sample"
  "verdmell/utils"
)
//
var env *environment.Environment

//#
//#
//# CheckSystem struct
//# The struct for CheckSystem has all Check, Checkgroup and samples information
type CheckSystem struct{
  // Map to storage the checks
  Ck *Checks  `json:"checks"`
  // Map to storage the checkgroups
  Cg *Checkgroups  `json:"checkgroups"`
  // Map to storage the samples
  Cs *sample.SampleSystem  `json:"samples"`
  // Timestamp
  Timestamp int64 `json:"timestamp"`
  // Service Channel
  outputSampleChan chan *sample.CheckSample `json: "-"`
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

  // get all check from Checks
  checkNames := ck.ListCheckNames()
  // set all checks list to the environment
  env.SetChecks(checkNames)

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
  // Initialize the first timestamp to 0
  cks.SetTimestamp(0)

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
//# SetTimestamp: attribute from CheckSystem
func (c *CheckSystem) SetTimestamp(t int64) {
  env.Output.WriteChDebug("(CheckSystem::SetTimestamp)")
  c.Timestamp = t
}
//
//# SetOutputSampleChan: methods sets the inputSampleChan's value
func (c *CheckSystem) SetOutputSampleChan(o chan *sample.CheckSample) {
  c.outputSampleChan = o
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
//# GetTimestamp: attribute from CheckSystem
func (c *CheckSystem) GetTimestamp() int64 {
  env.Output.WriteChDebug("(CheckSystem::GetTimestamp)")
  return c.Timestamp
}
//
//# GetOutputSampleChan: methods sets the inputSampleChan's value
func (c *CheckSystem) GetOutputSampleChan() chan *sample.CheckSample {
  return c.outputSampleChan
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# StartCheckSystem: will determine which kind of check has been required by user and start the checks
func (c *CheckSystem) StartCheckSystem(i interface{}) error {
  env.Output.WriteChDebug("(CheckSystem::StartCheckSystem)")

  endChan := make(chan bool)
  defer close(endChan)
  errChan := make(chan error)
  defer close(errChan)

  // the next will be only used during ec or eg
  // check will contain the Check configurations
  check := new(Checks)
  // checks will contain all the CheckObject definition
  checks := make(map[string]*CheckObject)
  
  //Increase the timestamp
  c.SetTimestamp(c.GetTimestamp()+1)

  switch req := i.(type){
  case *CheckObject:
    env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Starting the check '"+req.String()+"'")
    //add the check to be executed
    checks[req.GetName()] = req
    //add the check dependencies
    for _,dependency := range req.GetDepend(){
      if err, checkObj := c.Ck.GetCheckObjectByName(dependency); err != nil {
        return err
      } else {
        if _,exist := checks[dependency]; !exist{
          checks[dependency] = checkObj
        }
      }
    }
    check.SetCheck(checks)
    // run a goroutine for each checkObject and write the result to the channel
    go func() {
      // startCheckTaskPools requiere the SAmple system to sent sample to it and OutputSampleChan to send samples to ServiceSystem
      if err := check.StartCheckTaskPools(c.GetSampleSystem(),c.GetOutputSampleChan(),c.GetTimestamp()); err != nil {
        errChan <- err
      }
      endChan <- true
    }()

    select{
    case <-endChan:
      env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) All Pools Finished")
    case err := <-errChan:
      return err
    }

  case []string:
    env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Running a Checkgroup")

    for _,checkname := range req {
      env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) Preparing the check '"+checkname+"'")
      cks := c.GetChecks()

      if err, checkObj := cks.GetCheckObjectByName(checkname); err != nil {
        return err
      } else {
        if _,exist := checks[checkname]; !exist{
          checks[checkname] = checkObj
        }
        //add the check dependencies
        for _,dependency := range checkObj.GetDepend(){
          if err, checkObjdependency := c.Ck.GetCheckObjectByName(dependency); err != nil {
            return err
          } else {
            if _,exist := checks[dependency]; !exist{
              checks[dependency] = checkObjdependency
            }
          }
        }
      }
    }
    check.SetCheck(checks)
    // run a goroutine for each checkObject and write the result to the channel
    go func() {
      // startCheckTaskPools requiere the SAmple system to sent sample to it and OutputSampleChan to send samples to ServiceSystem
      if err := check.StartCheckTaskPools(c.GetSampleSystem(),c.GetOutputSampleChan(),c.GetTimestamp()); err != nil{
        errChan <- err
      }
      endChan <- true
    }()

    select{
    case <-endChan:
      env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) All Pools Finished")
    case err := <-errChan:
      return err
    }
    
  default:
    checks :=  c.GetChecks()
    // startCheckTaskPools requiere the SAmple system to sent sample to it and OutputSampleChan to send samples to ServiceSystem
    if err := checks.StartCheckTaskPools(c.GetSampleSystem(),c.GetOutputSampleChan(),c.GetTimestamp()); err != nil{
      return err
    }
    env.Output.WriteChDebug("(CheckSystem::StartCheckSystem) All Pools Finished")
  }

  return nil
}
//
//# InitCheckRunningQueues: prepares each checkobject to be run
func (c *CheckSystem) InitCheckRunningQueues() error {
  cs := c.GetChecks()

  for _,obj := range cs.GetCheck() {
      go func(checkObj *CheckObject) {
        env.Output.WriteChDebug("(CheckSystem::InitCheckRunningQueues) CheckQueue for '"+checkObj.GetName()+"'") 
        checkObj.StartQueue()
      }(obj)
  }
  return nil
}
//
//# AddCheck: method add a new check to the Checks struct
func (c *CheckSystem) AddCheck(obj *CheckObject) error {
  if err := c.Ck.AddCheck(obj); err != nil {
    return err
  }
  return nil
}
//
//# ListCheckNames: returns an array with the check namess defined on Checks object 
func (c *CheckSystem) ListCheckNames() []string {
  return c.Ck.ListCheckNames()
}
//
//# IsDefined: return if a check is defined
func (c *CheckSystem) IsDefined(name string) bool {
  return c.Ck.IsDefined(name)
}
//
//# GetCheckObjectByName: returns a check object gived a name
func (c *CheckSystem) GetCheckObjectByName(checkname string) (error, *CheckObject) {
  return c.Ck.GetCheckObjectByName(checkname)
}
//
//# GetCheckgroupByName: returns a check object gived a name
func (c *CheckSystem) GetCheckgroupByName(checkgroupname string) (error, []string) {
  return c.Cg.GetCheckgroupByName(checkgroupname)
}
//
//# GetAllChecks: return all checks
func (c *CheckSystem) GetAllChecks() string {
  env.Output.WriteChDebug("(CheckSystem::GetAllChecks)")
  // Get Checks attribute from CheckSystem
  checks := c.GetChecks()

  return checks.String()
}
//
//# GetCheck: return a checks
func (c *CheckSystem) GetCheck(check string) string {
  env.Output.WriteChDebug("(CheckSystem::GetCheck)")
  // Get Checks attribute from CheckSystem
  cks := c.GetChecks()
  // Get Check map from Checks
  ck := cks.GetCheck()
  // Get CheckObject
  obj := ck[check]
  return obj.String()
}
//
//# GetAllCheckgroups: return all checks
func (c *CheckSystem) GetAllCheckgroups() string {
  env.Output.WriteChDebug("(CheckSystem::GetAllCheckgroups)")
  // Get Checkgroups attribute from CheckSystem
  groups := c.GetCheckgroups()

  return groups.String()
}
//
//# GetCheckgroup: return a checks
func (c *CheckSystem) GetCheckgroup(group string) string {
  env.Output.WriteChDebug("(CheckSystem::GetCheckgroup)")
  // Get Checkgroupss attribute from CheckSystem
  cgs := c.GetCheckgroups()
  // Get Check map from Checks
  cg := cgs.GetCheckgroup()
  // Get CheckObject
  g := cg[group]
  return utils.ObjectToJsonString(g)
}

//
//# GetAllSamples: return the status of all checks
func (c *CheckSystem) GetAllSamples() string {
  env.Output.WriteChDebug("(CheckSystem::GetAllSamples)")
  // Get Samplesystem attribute from CheckSystem
  samplesystem := c.GetSampleSystem()

  return samplesystem.String()
}

//
//# GetSampleForCheck: return the status of all checks
func (c *CheckSystem) GetSampleForCheck(check string) string {
  env.Output.WriteChDebug("(CheckSystem::GetSampleForCheck)")
  
  samplesystem := c.GetSampleSystem()
  _,s := samplesystem.GetSample(check)

  return s.String()
}


//#
//# Common methods
//#---------------------------------------------------------------------

//# String: convert a Checks object to string
func (c *CheckSystem) String() string {
  return utils.ObjectToJsonString(c)
}

//#######################################################################################################