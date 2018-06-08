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
  "strconv"
  "verdmell/environment"
  "verdmell/sample"
  "verdmell/utils"

  "github.com/apenella/messageOutput"
)
// variable to set environment information
var env *environment.Environment
// manage messages/logger
var log *message.Message

//#
//#
//# CheckEngine struct
//# The struct for CheckEngine has all Check, Checkgroup and samples information
type CheckEngine struct{
  // Map to storage the checks
  Cks *Checks  `json:"checks"`
  // Map to storage the checkgroups
  Groups *Checkgroups  `json:"groups"`
  // Service Channel
  subscriptions map[chan interface{}] string `json: "-"`
}
//
//# NewCheckEngine: return a CheckEngine instance to be run
func NewCheckEngine(e *environment.Environment) (error, *CheckEngine){
  e.Output.WriteChDebug("(CheckEngine::NewCheckEngine) Create new engine instance")

  // get the environment attributes
  env = e
  // validate the sampleEngine status
  if sam := env.GetSampleEngine(); sam == nil {
    return errors.New("(CheckEngine::NewCheckEngine) SampleEngine have to be initialized before ChecksEngine's load"),nil
  }

  eng := new(CheckEngine)
  var err error

  // Get defined checks
  // validate checks and set the checks into check system
  cks := RetrieveChecks(env.Config.Checks.Folder)
  if err = cks.ValidateChecks(nil); err == nil {
    eng.SetChecks(cks)
    //Init the running queues to proceed the executions
    eng.InitCheckRunningQueues()
  } else {
    return err, nil
  }

  // Get defined checks groups
  // validate checks and set the checks into check system
  groups := RetrieveCheckgroups(env.Config.Checks.Folder)
  if err := groups.ValidateCheckgroups(cks); err == nil {
    eng.SetCheckgroups(groups)
  } else {
    return err, nil
  }

  // Initialize the Subscriptions
  eng.subscriptions = make(map[chan interface{}] string)

  // Set the environment's check engine
  env.SetCheckEngine(eng)
  // .WriteChInfo("(CheckEngine::NewCheckEngine) Hi! I'm your new check engine instance")

	return err, eng
}

//
// Interface Engine requirements

// Init
func (c *CheckEngine) Init() error {
  if log == nil {
    log = message.New(message.INFO,nil,0)
  }

  return nil
}
// Run
func (c *CheckEngine) Run() error {
  return nil
}
func (c *CheckEngine) Stop() error { return nil }
func (c *CheckEngine) Status() int { return 0 }

func (c *CheckEngine)	GetID() uint { return uint(0) }
func (c *CheckEngine)	GetName() string { return "" }
func (c *CheckEngine) GetDependencies() []uint { return nil }
func (c *CheckEngine) GetInputChannel() chan interface{} { return nil }
func (c *CheckEngine) GetStatus() uint { return uint(0) }
func (c *CheckEngine) SetStatus(s uint) {}

//
// Getters and Setters

// SetMsg
func SetMsg(l *message.Message) {
  log = l
}

//
//# SetChecks: attribute from CheckEngine
func (eng *CheckEngine) SetChecks(cks *Checks) {
  // env.Output.WriteChDebug("(CheckEngine::SetChecks) Set value '"+cks.String()+"'")
  log.Debug("(CheckEngine::SetChecks) Set value '"+cks.String()+"'")
  eng.Cks = cks
}
//
//# SetCheckgroups: attribute from CheckEngine
func (eng *CheckEngine) SetCheckgroups(groups *Checkgroups) {
  // env.Output.WriteChDebug("(CheckEngine::SetCheckgroups) Set value '"+groups.String()+"'")
  log.Debug("(CheckEngine::SetCheckgroups) Set value '"+groups.String()+"'")
  eng.Groups = groups
}
//
//# SetSubscriptions: method sets the channels to write samples
func (eng *CheckEngine) SetSubscriptions(o map[chan interface{}] string) {
  // env.Output.WriteChDebug("(CheckEngine::SetSubscriptions) Set value")
  log.Debug("(CheckEngine::SetSubscriptions) Set value")
  eng.subscriptions = o
}
//
//# Getchecks: attribute from CheckEngine
func (eng *CheckEngine) GetChecks() *Checks{
  // env.Output.WriteChDebug("(CheckEngine::GetChecks) Get value")
  log.Debug("(CheckEngine::GetChecks) Get value")
  return eng.Cks
}
//
//# Getcheckgroups: attribute from CheckEngine
func (eng *CheckEngine) GetCheckgroups() *Checkgroups{
  // env.Output.WriteChDebug("(CheckEngine::GetCheckgroups) Get value")
  log.Debug("(CheckEngine::GetCheckgroups) Get value")
  return eng.Groups
}
//
//# GetSubscriptions: methods return the channels to write samples
func (eng *CheckEngine) GetSubscriptions() map[chan interface{}] string {
  // env.Output.WriteChDebug("(CheckEngine::GetSubscriptions)")
  log.Debug("(CheckEngine::GetSubscriptions)")
  return eng.subscriptions
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# SayHi:
func (eng *CheckEngine) SayHi() {
  // env.Output.WriteChInfo("(CheckEngine::SayHi) Hi! I'm your new check engine instance")
  log.Info("(CheckEngine::SayHi) Hi! I'm your new check engine instance")
}
//
//# AddOutputSampleChan:
func (eng *CheckEngine) Subscribe(o chan interface{}, desc string) error {
  // env.Output.WriteChDebug("(CheckEngine::Subscribe)")
  log.Debug("(CheckEngine::Subscribe) ",desc)

  channels := eng.GetSubscriptions()
  if _, exist := channels[o]; !exist {
    channels[o] = desc
  } else {
    return errors.New("(CheckEngine::Subscribe) You are trying to add an existing channel")
  }

  return nil
}

//
//# InitCheckRunningQueues: prepares each checkobject to be run
func (eng *CheckEngine) InitCheckRunningQueues() error {
  // env.Output.WriteChDebug("(CheckEngine::InitCheckRunningQueues)")
  log.Debug("(CheckEngine::InitCheckRunningQueues)")
  cks := eng.GetChecks()

  for _,obj := range cks.GetCheck() {
      go func(checkObj *CheckObject) {
        // env.Output.WriteChDebug("(CheckEngine::InitCheckRunningQueues) CheckQueue for '"+checkObj.GetName()+"'")
        log.Debug("(CheckEngine::InitCheckRunningQueues) CheckQueue for '"+checkObj.GetName()+"'")
        checkObj.StartQueue()
      }(obj)
  }
  return nil
}
//
//# Start: will determine which kind of check has been required by user and start the checks
func (eng *CheckEngine) Start(i interface{}) error {
  // env.Output.WriteChDebug("(CheckEngine::Start)")
  log.Debug("(CheckEngine::Start)")
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
    // env.Output.WriteChDebug("(CheckEngine::Start) Starting the check '"+req.String()+"'")
    log.Debug("(CheckEngine::Start) Starting the check '"+req.String()+"'")
    //add the check to be executed
    checks[req.GetName()] = req
    //add the check dependencies
    for _,dependency := range req.GetDepend(){
      if err, checkObj := eng.Cks.GetCheckObjectByName(dependency); err != nil {
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
      // env.Output.WriteChDebug("(CheckEngine::Start) All Pools Finished")
      log.Debug("(CheckEngine::Start) All Pools Finished")
    case err := <-errChan:
      return err
    }

  case []string:
    // env.Output.WriteChDebug("(CheckEngine::Start) Running a Checkgroup")
    log.Debug("(CheckEngine::Start) Running a Checkgroup")
    for _,checkname := range req {
      // env.Output.WriteChDebug("(CheckEngine::Start) Preparing the check '"+checkname+"'")
      cks := eng.GetChecks()

      if err, checkObj := cks.GetCheckObjectByName(checkname); err != nil {
        return err
      } else {
        if _,exist := checks[checkname]; !exist{
          checks[checkname] = checkObj
        }
        //add the check dependencies
        for _,dependency := range checkObj.GetDepend(){
          if err, checkObjdependency := eng.Cks.GetCheckObjectByName(dependency); err != nil {
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
      // env.Output.WriteChDebug("(CheckEngine::Start) All Pools Finished")
      log.Debug("(CheckEngine::Start) All Pools Finished")
    case err := <-errChan:
      return err
    }

  default:
    checks :=  eng.GetChecks()
    // startCheckTaskPools requiere the Sample system to sent sample to it and OutputSampleChan to send samples to ServiceSystem
    if err := checks.StartCheckTaskPools(); err != nil{
      return err
    }
    // env.Output.WriteChDebug("(CheckEngine::Start) All Pools Finished")
    log.Debug("(CheckEngine::Start) All Pools Finished")
  }

  return nil
}
//
//# sendSamples: method that send samples to other engines
func (eng *CheckEngine) sendSample(s *sample.CheckSample) error {
  // env.Output.WriteChDebug("(CheckEngine::sendSample)["+strconv.Itoa(int(s.GetTimestamp()))+"] Send sample for '"+s.GetCheck()+"' check with exit '"+strconv.Itoa(s.GetExit())+"'")
  log.Debug("(CheckEngine::sendSample)["+strconv.Itoa(int(s.GetTimestamp()))+"] Send sample for '"+s.GetCheck()+"' check with exit '"+strconv.Itoa(s.GetExit())+"'")
  sampleEngine := env.GetSampleEngine().(*sample.SampleEngine)

  // send samples to ServiceEngine
  // GetSample return an error if no samples has been add before for that check
  if err, sam := sampleEngine.GetSample(s.GetCheck()); err != nil {
    // sending sample to service using the output channel
    eng.notify(s)
  } else {
    // if a sample for that exist
    // the sample will not send to service system unless it has modified it exit status
    if sam.GetTimestamp() < s.GetTimestamp() {
      // sending sample to service using the output channel
      eng.notify(s)
    }
  }

  return nil
}
//
//# writeToSubscriptions: function write samples to defined samples
func (eng *CheckEngine) notify(s *sample.CheckSample) {
  // env.Output.WriteChDebug("(CheckEngine::notify)")
  log.Debug("(CheckEngine::notify)")
  for o, desc := range eng.GetSubscriptions(){
    // env.Output.WriteChDebug("(CheckEngine::notify) ["+strconv.Itoa(int(s.GetTimestamp()))+"] Notify sample '"+s.GetCheck()+"' with exit '"+strconv.Itoa(s.GetExit())+"' on channel '"+desc+"'")
    log.Debug("(CheckEngine::notify) ["+strconv.Itoa(int(s.GetTimestamp()))+"] Notify sample '"+s.GetCheck()+"' with exit '"+strconv.Itoa(s.GetExit())+"' on channel '"+desc+"'")
    o <- s
  }
}

//
//# AddCheck: method add a new check to the Checks struct
func (eng *CheckEngine) AddCheck(obj *CheckObject) error {
  if err := eng.Cks.AddCheck(obj); err != nil {
    return err
  }
  return nil
}
//
//# ListCheckNames: returns an array with the check namess defined on Checks object
func (eng *CheckEngine) ListCheckNames() []string {
  return eng.Cks.ListCheckNames()
}
//
//# IsDefined: return if a check is defined
func (eng *CheckEngine) IsDefined(name string) bool {
  return eng.Cks.IsDefined(name)
}
//
//# GetCheckObjectByName: returns a check object gived a name
func (eng *CheckEngine) GetCheckObjectByName(checkname string) (error, *CheckObject) {
  return eng.Cks.GetCheckObjectByName(checkname)
}
//
//# GetCheckgroupByName: returns a check object gived a name
func (eng *CheckEngine) GetCheckgroupByName(checkgroupname string) (error, []string) {
  return eng.Groups.GetCheckgroupByName(checkgroupname)
}
//
//# GetAllChecks: return all checks
func (eng *CheckEngine) GetAllChecks() (error,[]byte) {
  // env.Output.WriteChDebug("(CheckEngine::GetAllChecks)")
  var checks *Checks

  if checks = eng.GetChecks(); checks == nil {
    msg := "(CheckEngine::GetAllChecks) There are no checks defined."
    // env.Output.WriteChDebug(msg)
    log.Debug(msg)
    return errors.New(msg), nil
  }

  return utils.ObjectToJsonByte(checks)
}
//
//# GetCheck: return a checks
func (eng *CheckEngine) GetCheck(name string) (error,[]byte) {
  // env.Output.WriteChDebug("(CheckEngine::GetCheck)")
  var checks *Checks
  var check map[string] *CheckObject
  var obj *CheckObject
  var exist bool

  // Get Checks attribute from CheckEngine
  if checks = eng.GetChecks(); checks == nil {
    msg := "(CheckEngine::GetCheck) There are no checks defined."
    // env.Output.WriteChDebug(msg)
    log.Debug(msg)
    return errors.New(msg), nil
  }
  // Get Check map from Checks
  if check = checks.GetCheck(); checks == nil {
    msg := "(CheckEngine::GetCheck) There are no checks defined."
    // env.Output.WriteChDebug(msg)
    log.Debug(msg)
    return errors.New(msg), nil
  }
  // Get CheckObject from the check's map
  if obj,exist = check[name]; !exist {
    msg := "(ServiceEngine::GetCheck) The check '"+name+"' is not defined."
    // env.Output.WriteChDebug(msg)
    log.Debug(msg)
    return errors.New(msg), nil
  }

  return utils.ObjectToJsonByte(obj)
}
//
//# GetAllCheckgroups: return all checks
func (eng *CheckEngine) GetAllCheckgroups() (error,[]byte) {
  // env.Output.WriteChDebug("(CheckEngine::GetAllCheckgroups)")
  var groups *Checkgroups

  if groups = eng.GetCheckgroups(); groups == nil {
    msg := "(CheckEngine::GetAllCheckgroups) There are no check groups defined."
    // env.Output.WriteChDebug(msg)
    log.Debug(msg)
    return errors.New(msg), nil
  }

  return utils.ObjectToJsonByte(groups)
}
//
//# GetCheckgroup: return a checks
func (eng *CheckEngine) GetCheckgroup(name string) (error,[]byte) {
  // env.Output.WriteChDebug("(CheckEngine::GetCheckgroup)")
  var groups *Checkgroups
  var group map[string] []string
  var obj []string
  var exist bool

  // Get Checkgroupss attribute from CheckEngine
  if groups = eng.GetCheckgroups(); groups == nil {
    msg := "(CheckEngine::GetCheckgroup) There are no check groups defined."
    // env.Output.WriteChDebug(msg)
    log.Debug(msg)
    return errors.New(msg), nil
  }
  // Get Check map from Checks
  if group = groups.GetCheckgroup(); group == nil {
    msg := "(CheckEngine::GetCheckgroup) There are no check groups defined."
    // env.Output.WriteChDebug(msg)
    log.Debug(msg)
    return errors.New(msg), nil
  }
  // Get Check group from check group's mpa
  if obj,exist = group[name]; !exist {
    msg := "(ServiceEngine::GetCheckgroup) The check group '"+name+"' is not defined."
    // env.Output.WriteChDebug(msg)
    log.Debug(msg)
    return errors.New(msg), nil
  }

  return utils.ObjectToJsonByte(obj)
}

//#
//# Common methods
//#---------------------------------------------------------------------

//# String: convert a Checks object to string
func (eng *CheckEngine) String() string {
  if err, str := utils.ObjectToJsonString(eng); err != nil{
    return err.Error()
  } else{
    return str
  }
}

//#######################################################################################################
