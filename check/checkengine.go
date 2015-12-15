/*
Check Engine management

The package 'check' is used by verdmell to manage the monitoring checks defined by user

-Checks
-CheckObject
-Checkgroups

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
//# CheckEngine struct
//# The struct for CheckEngine has all Check, Checkgroup and samples information
type CheckEngine struct{
  // Map to storage the checks
  Ck *Checks  `json:"checks"`
  // Map to storage the checkgroups
  Cg *Checkgroups  `json:"checkgroups"`
  // Map to storage the samples
  Cs *sample.SampleEngine  `json:"samples"`
  // Timestamp
  Timestamp int64 `json:"timestamp"`
  // Service Channel
  outputSampleChan chan *sample.CheckSample `json: "-"`
}
//
//# NewCheckEngine: return a CheckEngine instance to be run
func NewCheckEngine(e *environment.Environment) (error, *CheckEngine){
  e.Output.WriteChDebug("(CheckEngine::NewCheckEngine)")
  cks := new(CheckEngine)
  ss := new(sample.SampleEngine)
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
  if err, ss = sample.NewSampleEngine(env); err == nil {
   cks.SetSampleEngine(ss)
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
//# SetChecks: attribute from CheckEngine
func (c *CheckEngine) SetChecks(ck *Checks) {
  env.Output.WriteChDebug("(CheckEngine::SetChecks) Set value '"+ck.String()+"'")
  c.Ck = ck
}
//
//# SetCheckgroups: attribute from CheckEngine
func (c *CheckEngine) SetCheckgroups(cg *Checkgroups) {
  env.Output.WriteChDebug("(CheckEngine::SetCheckgroups) Set value '"+cg.String()+"'")
  c.Cg = cg
}
//
//# SetChecksamplesmap: attribute from CheckEngine
func (c *CheckEngine) SetSampleEngine(cs *sample.SampleEngine) {
  env.Output.WriteChDebug("(CheckEngine::SetSampleEngine)")
  c.Cs = cs
}
//
//# SetTimestamp: attribute from CheckEngine
func (c *CheckEngine) SetTimestamp(t int64) {
  env.Output.WriteChDebug("(CheckEngine::SetTimestamp)")
  c.Timestamp = t
}
//
//# SetOutputSampleChan: methods sets the inputSampleChan's value
func (c *CheckEngine) SetOutputSampleChan(o chan *sample.CheckSample) {
  c.outputSampleChan = o
}

//
//# Getchecks: attribute from CheckEngine
func (c *CheckEngine) GetChecks() *Checks{
  env.Output.WriteChDebug("(CheckEngine::GetChecks) Get value")
  return c.Ck
}
//
//# Getcheckgroups: attribute from CheckEngine
func (c *CheckEngine) GetCheckgroups() *Checkgroups{
  env.Output.WriteChDebug("(CheckEngine::GetCheckgroups) Get value")
  return c.Cg
}
//
//# GetChecksamplesmap: attribute from CheckEngine
func (c *CheckEngine) GetSampleEngine() *sample.SampleEngine{
  env.Output.WriteChDebug("(CheckEngine::GetSampleEngine)")
  return c.Cs
}
//
//# GetTimestamp: attribute from CheckEngine
func (c *CheckEngine) GetTimestamp() int64 {
  env.Output.WriteChDebug("(CheckEngine::GetTimestamp)")
  return c.Timestamp
}
//
//# GetOutputSampleChan: methods sets the inputSampleChan's value
func (c *CheckEngine) GetOutputSampleChan() chan *sample.CheckSample {
  return c.outputSampleChan
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# StartCheckEngine: will determine which kind of check has been required by user and start the checks
func (c *CheckEngine) StartCheckEngine(i interface{}) error {
  env.Output.WriteChDebug("(CheckEngine::StartCheckEngine)")

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
    env.Output.WriteChDebug("(CheckEngine::StartCheckEngine) Starting the check '"+req.String()+"'")
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
      if err := check.StartCheckTaskPools(c.GetSampleEngine(),c.GetOutputSampleChan(),c.GetTimestamp()); err != nil {
        errChan <- err
      }
      endChan <- true
    }()

    select{
    case <-endChan:
      env.Output.WriteChDebug("(CheckEngine::StartCheckEngine) All Pools Finished")
    case err := <-errChan:
      return err
    }

  case []string:
    env.Output.WriteChDebug("(CheckEngine::StartCheckEngine) Running a Checkgroup")

    for _,checkname := range req {
      env.Output.WriteChDebug("(CheckEngine::StartCheckEngine) Preparing the check '"+checkname+"'")
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
      if err := check.StartCheckTaskPools(c.GetSampleEngine(),c.GetOutputSampleChan(),c.GetTimestamp()); err != nil{
        errChan <- err
      }
      endChan <- true
    }()

    select{
    case <-endChan:
      env.Output.WriteChDebug("(CheckEngine::StartCheckEngine) All Pools Finished")
    case err := <-errChan:
      return err
    }
    
  default:
    checks :=  c.GetChecks()
    // startCheckTaskPools requiere the SAmple system to sent sample to it and OutputSampleChan to send samples to ServiceSystem
    if err := checks.StartCheckTaskPools(c.GetSampleEngine(),c.GetOutputSampleChan(),c.GetTimestamp()); err != nil{
      return err
    }
    env.Output.WriteChDebug("(CheckEngine::StartCheckEngine) All Pools Finished")
  }

  return nil
}
//
//# InitCheckRunningQueues: prepares each checkobject to be run
func (c *CheckEngine) InitCheckRunningQueues() error {
  cs := c.GetChecks()

  for _,obj := range cs.GetCheck() {
      go func(checkObj *CheckObject) {
        env.Output.WriteChDebug("(CheckEngine::InitCheckRunningQueues) CheckQueue for '"+checkObj.GetName()+"'") 
        checkObj.StartQueue()
      }(obj)
  }
  return nil
}
//
//# AddCheck: method add a new check to the Checks struct
func (c *CheckEngine) AddCheck(obj *CheckObject) error {
  if err := c.Ck.AddCheck(obj); err != nil {
    return err
  }
  return nil
}
//
//# ListCheckNames: returns an array with the check namess defined on Checks object 
func (c *CheckEngine) ListCheckNames() []string {
  return c.Ck.ListCheckNames()
}
//
//# IsDefined: return if a check is defined
func (c *CheckEngine) IsDefined(name string) bool {
  return c.Ck.IsDefined(name)
}
//
//# GetCheckObjectByName: returns a check object gived a name
func (c *CheckEngine) GetCheckObjectByName(checkname string) (error, *CheckObject) {
  return c.Ck.GetCheckObjectByName(checkname)
}
//
//# GetCheckgroupByName: returns a check object gived a name
func (c *CheckEngine) GetCheckgroupByName(checkgroupname string) (error, []string) {
  return c.Cg.GetCheckgroupByName(checkgroupname)
}



//
//# GetAllChecks: return all checks
func (c *CheckEngine) GetAllChecks() []byte {
  env.Output.WriteChDebug("(CheckEngine::GetAllChecks)")
  return utils.ObjectToJsonByte(c.GetChecks())
}
//
//# GetCheck: return a checks
func (c *CheckEngine) GetCheck(check string) []byte {
  env.Output.WriteChDebug("(CheckEngine::GetCheck)")
  // Get Checks attribute from CheckEngine
  cks := c.GetChecks()
  // Get Check map from Checks
  ck := cks.GetCheck()
  return utils.ObjectToJsonByte(ck[check])
}
//
//# GetAllCheckgroups: return all checks
func (c *CheckEngine) GetAllCheckgroups() []byte {
  env.Output.WriteChDebug("(CheckEngine::GetAllCheckgroups)")
  return utils.ObjectToJsonByte(c.GetCheckgroups())
}
//
//# GetCheckgroup: return a checks
func (c *CheckEngine) GetCheckgroup(group string) []byte {
  env.Output.WriteChDebug("(CheckEngine::GetCheckgroup)")
  // Get Checkgroupss attribute from CheckEngine
  cgs := c.GetCheckgroups()
  // Get Check map from Checks
  cg := cgs.GetCheckgroup()
  return utils.ObjectToJsonByte(cg[group])
}

//
//# GetAllSamples: return the status of all checks
func (c *CheckEngine) GetAllSamples() []byte {
  env.Output.WriteChDebug("(CheckEngine::GetAllSamples)")
  return utils.ObjectToJsonByte(c.GetSampleEngine()) 
}

//
//# GetSampleForCheck: return the status of all checks
func (c *CheckEngine) GetSampleForCheck(check string) []byte {
  env.Output.WriteChDebug("(CheckEngine::GetSampleForCheck)")
  SampleEngine := c.GetSampleEngine()
  _,s := SampleEngine.GetSample(check)
  return utils.ObjectToJsonByte(s)
}


//#
//# Common methods
//#---------------------------------------------------------------------

//# String: convert a Checks object to string
func (c *CheckEngine) String() string {
  return utils.ObjectToJsonString(c)
}

//#######################################################################################################