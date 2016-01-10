/*
Check Engine management

The package 'check' is used by verdmell to manage the monitoring checks defined by user

-Checks
-CheckObject
-Checkgroups

*/
package check

import (
  "errors"
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
  // Service Channel
  outputSampleChan chan *sample.CheckSample `json: "-"`
}
//
//# NewCheckEngine: return a CheckEngine instance to be run
func NewCheckEngine(e *environment.Environment) (error, *CheckEngine){
  e.Output.WriteChDebug("(CheckEngine::NewCheckEngine)")

  // get the environment attributes
  env = e
  // validate the sampleEngine status
  if sam := env.GetSampleEngine(); sam == nil {
    return errors.New("(CheckEngine::NewCheckEngine) SampleEngine have to be initialized before ChecksEngine's load"),nil
  }

  cks := new(CheckEngine)
  var err error

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

  // Set the environment's check engine
  env.SetCheckEngine(cks)

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
//# GetOutputSampleChan: methods sets the inputSampleChan's value
func (c *CheckEngine) GetOutputSampleChan() chan *sample.CheckSample {
  return c.outputSampleChan
}

//#
//# Specific methods
//#---------------------------------------------------------------------

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
      // startCheckTaskPools requiere the Sample system to sent sample to it and OutputSampleChan to send samples to ServiceSystem
      if err := check.StartCheckTaskPools(); err != nil {
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
      if err := check.StartCheckTaskPools(); err != nil{
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
    // startCheckTaskPools requiere the Sample system to sent sample to it and OutputSampleChan to send samples to ServiceSystem
    if err := checks.StartCheckTaskPools(); err != nil{
      return err
    }
    env.Output.WriteChDebug("(CheckEngine::StartCheckEngine) All Pools Finished")
  }

  return nil
}
//
//# sendSamples: method that send samples to other engines
func (c *CheckEngine) sendSample(s *sample.CheckSample) error {
  env.Output.WriteChDebug("(CheckEngine::sendSample) Send sample for '"+s.GetCheck()+"' check")
  sampleEngine := env.GetSampleEngine().(*sample.SampleEngine)

  // send samples to ServiceEngine
  // GetSample return an error if no samples has been add before for that check
  if err, sam := sampleEngine.GetSample(s.GetCheck()); err != nil {
    // sending sample to service using the output channel
    c.outputSampleChan <- s
  } else {
    // if a sample for that exist
    // the sample will not send to service system unless it has modified it exit status
    if sam.GetTimestamp() < s.GetTimestamp() {
      // sending sample to service using the output channel
      c.outputSampleChan <- s
    }
  }

  // Add samples to SampleEngine
  sampleEngine.AddSample(s)
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
func (c *CheckEngine) GetAllChecks() (error,[]byte) {
  env.Output.WriteChDebug("(CheckEngine::GetAllChecks)")
  var checks *Checks

  if checks = c.GetChecks(); checks == nil {
    msg := "(CheckEngine::GetAllChecks) There are no checks defined."
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  } 

  return nil,utils.ObjectToJsonByte(checks)
}
//
//# GetCheck: return a checks
func (c *CheckEngine) GetCheck(name string) (error,[]byte) {
  env.Output.WriteChDebug("(CheckEngine::GetCheck)")
  var checks *Checks
  var check map[string] *CheckObject
  var obj *CheckObject
  var exist bool

  // Get Checks attribute from CheckEngine
  if checks = c.GetChecks(); checks == nil {
    msg := "(CheckEngine::GetCheck) There are no checks defined."
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  }
  // Get Check map from Checks
  if check = checks.GetCheck(); checks == nil {
    msg := "(CheckEngine::GetCheck) There are no checks defined."
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  }
  // Get CheckObject from the check's map
  if obj,exist = check[name]; !exist {
    msg := "(ServiceEngine::GetCheck) The check '"+name+"' is not defined."
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  }

  return nil,utils.ObjectToJsonByte(obj)
}
//
//# GetAllCheckgroups: return all checks
func (c *CheckEngine) GetAllCheckgroups() (error,[]byte) {
  env.Output.WriteChDebug("(CheckEngine::GetAllCheckgroups)")
  var groups *Checkgroups

  if groups = c.GetCheckgroups(); groups == nil {
    msg := "(CheckEngine::GetAllCheckgroups) There are no check groups defined."
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  }

  return nil,utils.ObjectToJsonByte(groups)
}
//
//# GetCheckgroup: return a checks
func (c *CheckEngine) GetCheckgroup(name string) (error,[]byte) {
  env.Output.WriteChDebug("(CheckEngine::GetCheckgroup)")
  var groups *Checkgroups
  var group map[string] []string
  var obj []string
  var exist bool

  // Get Checkgroupss attribute from CheckEngine
  if groups = c.GetCheckgroups(); groups == nil {
    msg := "(CheckEngine::GetCheckgroup) There are no check groups defined."
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  }
  // Get Check map from Checks
  if group = groups.GetCheckgroup(); group == nil {
    msg := "(CheckEngine::GetCheckgroup) There are no check groups defined."
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  }
  // Get Check group from check group's mpa
  if obj,exist = group[name]; !exist {
    msg := "(ServiceEngine::GetCheckgroup) The check group '"+name+"' is not defined."
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  }
  
  return nil,utils.ObjectToJsonByte(obj)
}

//#
//# Common methods
//#---------------------------------------------------------------------

//# String: convert a Checks object to string
func (c *CheckEngine) String() string {
  if err, str := utils.ObjectToJsonString(c); err != nil{
    return err.Error()
  } else{
    return str
  }
}

//#######################################################################################################