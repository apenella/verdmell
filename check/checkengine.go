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
  "verdmell/configuration"
  "verdmell/environment"
  "verdmell/sample"
  "verdmell/utils"

  "github.com/apenella/messageOutput"
)

var env *environment.Environment

//#
//#
//# CheckEngine struct
//# The struct for CheckEngine has all Check, Checkgroup and samples information
type CheckEngine struct{
  // Map to storage the checks
  Cks *Checks  `json:"checks"`
  // Service Channel
  subscriptions map[chan interface{}] string `json: "-"`
  // variable to set configuration
  config *configuration.Configuration `json: "-"`
}
//
//# NewCheckEngine: return a CheckEngine instance to be run
func NewCheckEngine(c *configuration.Configuration) *CheckEngine {
  c.Log.Debug("(CheckEngine::NewCheckEngine) Create new engine instance")

  var err error

  eng := &CheckEngine{
    Cks: NewChecks(c),
    subscriptions: make(map[chan interface{}] string),
    config: c,
  }

	return eng
}

//
// Interface Engine requirements

// Init
func (eng *CheckEngine) Init() error {
  // initialize eng.config.Log.Er
  if eng.config.Log == nil {
    eng.config.Log= message.New(message.INFO,nil,0)
  }

  // Get defined checks
  // validate checks and set the checks into check system
  cks := RetrieveChecks(eng.config.Checks.Folder)
  if err := cks.ValidateChecks(nil); err == nil {
    eng.SetChecks(cks)
    //Init the running queues to proceed the executions
    eng.InitCheckRunningQueues()
  } else {
    return err
  }

  return nil
}

// Run
func (eng *CheckEngine) Run() error {
  return nil
}
func (eng *CheckEngine) Stop() error { return nil }
func (eng *CheckEngine) Status() int { return 0 }
func (eng *CheckEngine)	GetID() uint { return uint(0) }
func (eng *CheckEngine)	GetName() string { return "" }
func (eng *CheckEngine) GetDependencies() []uint { return nil }
func (eng *CheckEngine) GetInputChannel() chan interface{} { return nil }
func (eng *CheckEngine) GetStatus() uint { return uint(0) }
func (eng *CheckEngine) SetStatus(s uint) {}

//
// Getters and Setters

//
//# SetChecks: attribute from CheckEngine
func (eng *CheckEngine) SetChecks(cks *Checks) {
  eng.config.Log.Debug("(CheckEngine::SetChecks) Set value '"+cks.String()+"'")
  eng.Cks = cks
}

//
//# SetSubscriptions: method sets the channels to write samples
func (eng *CheckEngine) SetSubscriptions(o map[chan interface{}] string) {
  eng.config.Log.Debug("(CheckEngine::SetSubscriptions) Set value")
  eng.subscriptions = o
}
//
//# Getchecks: attribute from CheckEngine
func (eng *CheckEngine) GetChecks() *Checks{
  eng.config.Log.Debug("(CheckEngine::GetChecks) Get value")
  return eng.Cks
}
//
//# GetSubscriptions: methods return the channels to write samples
func (eng *CheckEngine) GetSubscriptions() map[chan interface{}] string {
  eng.config.Log.Debug("(CheckEngine::GetSubscriptions)")
  return eng.subscriptions
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# SayHi:
func (eng *CheckEngine) SayHi() {
  eng.config.Log.Info("(CheckEngine::SayHi) Hi! I'm your new check engine instance")
}
//
//# AddOutputSampleChan:
func (eng *CheckEngine) Subscribe(o chan interface{}, desc string) error {
  eng.config.Log.Debug("(CheckEngine::Subscribe) ",desc)

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
  eng.config.Log.Debug("(CheckEngine::InitCheckRunningQueues)")
  cks := eng.GetChecks()

  for _,obj := range cks.GetCheck() {
      go func(checkObj *CheckObject) {
        eng.config.Log.Debug("(CheckEngine::InitCheckRunningQueues) CheckQueue for '"+checkObj.GetName()+"'")
        checkObj.StartQueue()
      }(obj)
  }
  return nil
}
//
//# Start: will determine which kind of check has been required by user and start the checks
func (eng *CheckEngine) Start(i interface{}) error {
  eng.config.Log.Debug("(CheckEngine::Start)")
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
    eng.config.Log.Debug("(CheckEngine::Start) Starting the check '"+req.String()+"'")
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
      if err := check.StartCheckTaskPools(); err != nil {
        errChan <- err
      }
      endChan <- true
    }()

    select{
    case <-endChan:
      eng.config.Log.Debug("(CheckEngine::Start) All Pools Finished")
    case err := <-errChan:
      return err
    }

  case []string:
    eng.config.Log.Debug("(CheckEngine::Start) Running a Checkgroup")
    for _,checkname := range req {
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
      eng.config.Log.Debug("(CheckEngine::Start) All Pools Finished")
    case err := <-errChan:
      return err
    }

  default:
    checks :=  eng.GetChecks()
    if err := checks.StartCheckTaskPools(); err != nil{
      return err
    }
    eng.config.Log.Debug("(CheckEngine::Start) All Pools Finished")
  }

  return nil
}
//
//# sendSamples: method that send samples to other engines
func (eng *CheckEngine) sendSample(s *sample.CheckSample) error {
  eng.config.Log.Debug("(CheckEngine::sendSample)["+strconv.Itoa(int(s.GetTimestamp()))+"] Send sample for '"+s.GetCheck()+"' check with exit '"+strconv.Itoa(s.GetExit())+"'")
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
  eng.config.Log.Debug("(CheckEngine::notify)")
  for o, desc := range eng.GetSubscriptions(){
    eng.config.Log.Debug("(CheckEngine::notify) ["+strconv.Itoa(int(s.GetTimestamp()))+"] Notify sample '"+s.GetCheck()+"' with exit '"+strconv.Itoa(s.GetExit())+"' on channel '"+desc+"'")
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
//# GetAllChecks: return all checks
func (eng *CheckEngine) GetAllChecks() (error,[]byte) {
  var checks *Checks

  if checks = eng.GetChecks(); checks == nil {
    msg := "(CheckEngine::GetAllChecks) There are no checks defined."
    eng.config.Log.Debug(msg)
    return errors.New(msg), nil
  }

  return utils.ObjectToJsonByte(checks)
}
//
//# GetCheck: return a checks
func (eng *CheckEngine) GetCheck(name string) (error,[]byte) {
  var checks *Checks
  var check map[string] *CheckObject
  var obj *CheckObject
  var exist bool

  // Get Checks attribute from CheckEngine
  if checks = eng.GetChecks(); checks == nil {
    msg := "(CheckEngine::GetCheck) There are no checks defined."
    eng.config.Log.Debug(msg)
    return errors.New(msg), nil
  }
  // Get Check map from Checks
  if check = checks.GetCheck(); checks == nil {
    msg := "(CheckEngine::GetCheck) There are no checks defined."
    eng.config.Log.Debug(msg)
    return errors.New(msg), nil
  }
  // Get CheckObject from the check's map
  if obj,exist = check[name]; !exist {
    msg := "(ServiceEngine::GetCheck) The check '"+name+"' is not defined."
    eng.config.Log.Debug(msg)
    return errors.New(msg), nil
  }

  return utils.ObjectToJsonByte(obj)
}

//#
//# Common methods

//
//# String: convert a Checks object to string
func (eng *CheckEngine) String() string {
  if err, str := utils.ObjectToJsonString(eng); err != nil{
    return err.Error()
  } else{
    return str
  }
}
