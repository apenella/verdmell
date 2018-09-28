/*
Check Engine management

The package 'check' is used by verdmell to manage the monitoring checks defined by user

-Checks
-CheckObject
-Checkgroups

*/
package check

import (
  "os/exec"
  "syscall"
  "errors"
  "time"
  "strings"
  "strconv"
  "verdmell/sample"
  "verdmell/utils"
)
//
//
// CheckObject stuct:
// checkobject defines a check to be executed
type CheckObject struct{
  Name string `json:"name"`
  Description string `json:"description"`
  Command string `json:"command"`
  Depend []string `json:"depend"`
  ExpirationTime int `json:"expiration_time"`
  Interval int `json:"interval"`
  Custom interface{} `json:"-"`
  // Timestamp
  Timestamp int64 `json:"timestamp"`
  //Queues
  TaskQueue chan *CheckObject `json:"-"`
  //StatusChan chan int
  SampleChan chan *sample.CheckSample `json:"-"`
}

//
// Getters/Setters methods for Checks object

//
// SetName: method sets the Name value for the CheckObject object
func (c *CheckObject) SetName(n string) {
  c.Name = n
}
//
// SetDescription: method sets the Description value for the CheckObject object
func (c *CheckObject) SetDescription(d string) {
  c.Description = d
}
//
// SetCommand: method sets the Command value for the CheckObject object
func (c *CheckObject) SetCommand(cmd string) {
  c.Command = cmd
}
//
// SetDepend: method sets the Depend value for the CheckObject object
func (c *CheckObject) SetDepend(d []string) {
  c.Depend = d
}
//
// SetExpirationTime: method sets the ExpirationTime value for the CheckObject object
func (c *CheckObject) SetExpirationTime(t int) {
  c.ExpirationTime = t
}
//
// SetInterval: method sets the Interval value for the CheckObject object
func (c *CheckObject) SetInterval(i int) {
  c.Interval = i
}
//
// SetTimestamp: attribute from CheckObject
func (c *CheckObject) SetTimestamp(t int64) {
  env.Output.WriteChDebug("(CheckObject::SetTimestamp)")
  c.Timestamp = t
}
//
// SetTaskQueue: method sets the queue value for the CheckObject object
func (c *CheckObject) SetTaskQueue(q chan *CheckObject) {
  c.TaskQueue = q
}
//
// SetSampleChan: method sets the StatusChan value for the CheckObject object
func (c *CheckObject) SetSampleChan(sc chan *sample.CheckSample) {
  c.SampleChan = sc
}

//
// GetName: method returns the Name value for the CheckObject object
func (c *CheckObject) GetName() string {
  return c.Name
}
//
// GetDescription: method returns the Description value for the CheckObject object
func (c *CheckObject) GetDescription() string {
  return c.Description
}
//
// GetCommand: method returns the Command value for the CheckObject object
func (c *CheckObject) GetCommand() string {
  return c.Command
}
//
// GetDepend: method returns the Depend value for the CheckObject object
func (c *CheckObject) GetDepend() []string {
  return c.Depend
}
//
// GetExpirationTime: method returns the ExpirationTime value for the CheckObject object
func (c *CheckObject) GetExpirationTime() int{
  return c.ExpirationTime
}
//
// GetInterval: method returns the Interval value for the CheckObject object
func (c *CheckObject) GetInterval() int{
  return c.Interval
}
//
// GetTimestamp: attribute from CheckObject
func (c *CheckObject) GetTimestamp() int64 {
  return c.Timestamp
}
//
// GetTaskQueue: method returns the TaskQueue value for the CheckObject object
func (c *CheckObject) GetTaskQueue() chan *CheckObject{
  return c.TaskQueue
}
//
// GetSampleChan: method returns the StatusChan value for the CheckObject object
func (c *CheckObject) GetSampleChan() chan *sample.CheckSample{
  return c.SampleChan
}

//
// ValidateCheckObject: ensures that the CheckObject has all the required data set. The method returns an error object once a definition method is found
func (c *CheckObject) ValidateCheckObject() error {

  env.Output.WriteChDebug("(CheckObject::ValidateCheckObject) Check '"+c.GetName()+"'")
  if c.Command == "" {
    err := errors.New("(CheckObject::ValidateCheckObject) Check '"+c.GetName()+"' requires a Command")
    return err
  }
  if c.GetExpirationTime() < 0 {
    err := errors.New("(CheckObject::ValidateCheckObject) Check '"+c.GetName()+"' has an invalid expiration time")
    return err
  }
  if c.GetInterval() < 0 {
    err := errors.New("(CheckObject::ValidateCheckObject) Check '"+c.GetName()+"' has an invalid interval")
    return err
  }

  env.Output.WriteChDebug("(CheckObject::ValidateCheckObject) Check '"+c.GetName()+"' is correct")
  return nil
}
//
// StartQueue: method starts a queue for receive check
func (c *CheckObject) StartQueue(){
  env.Output.WriteChDebug("(CheckObject::StartQueue) Starting queue for check '"+c.GetName()+"'")
  var err error
  expired := make(chan bool)
  result := -1
  queue := c.TaskQueue
  sample := new(sample.CheckSample)

  defer close(queue)

  //function to clean up the result to enforce that the check is not started while the sample is already valid
  sampleExpiration := func() {
      env.Output.WriteChDebug("(CheckObject::StartQueue::sampleExpiration) Countdown for "+c.GetName()+"'s sample")
      timeout := time.After(time.Duration(c.GetExpirationTime()) * time.Second)
      for{
        select{
        case <-timeout:
          expired <- true
        }
      }
  }

  scheduleCheckTask := func () {
    env.Output.WriteChDebug("(CheckObject::StartQueue::scheduleCheckTask) Scheduling a new task for '"+c.GetName()+"'. It will be launched "+strconv.Itoa(c.GetInterval())+"s latter.")
    timeout := time.After(time.Duration(c.GetInterval()) * time.Second)
    for{
      select{
      case <-timeout:
        c.SetTimestamp(c.GetTimestamp()+1)
        checkEngine := env.GetCheckEngine().(*CheckEngine)
        checkEngine.Start(c)
      }
    }
  }

  // waiting for task to be queued by EnqueueCheckObject
  for{
    select{
    case checkObj := <-queue:
      if result >= 0 {
        env.Output.WriteChDebug("(CheckObject::StartQueue) ObjectTask alive and won't be started again. Check '"+checkObj.GetName()+"' already has a sample")
      } else {
        env.Output.WriteChDebug("(CheckObject::StartQueue) ObjectTask started. Check '"+checkObj.GetName()+"' has no sample")
        if err,sample = checkObj.StartCheckObjectTask(); err != nil {
          env.Output.WriteChWarn("(CheckObject::StartQueue) Task for '"+checkObj.GetName()+"' has not finished properly")
        }
        result = sample.GetExit()
        env.Output.WriteChDebug("(CheckObject::StartQueue) ObjectTask finished. Exit code for '"+checkObj.GetName()+"' is '"+strconv.Itoa(result)+"'")

        go sampleExpiration()
        go scheduleCheckTask()
      }
      //Send sample to check object sampleChan.
      checkObj.SampleChan <- sample
    case <-expired:
      env.Output.WriteChDebug("(CheckObject::StartQueue) Sample for "+c.GetName()+" has expired")
      result = -1
    }
  }
}
//
// EnqueueCheckObject: enqueu a CheckObject to be run
func (c *CheckObject) EnqueueCheckObject() (error) {
  env.Output.WriteChDebug("(CheckObject::EnqueueCheckObject) Enqueing CheckObject '"+c.GetName()+"'")
  c.TaskQueue <- c

  return nil
}
//
// StartCheckObjectTask: executes the command defined on check an return the result
func (c *CheckObject) StartCheckObjectTask() (error,  *sample.CheckSample) {
  env.Output.WriteChDebug("(CheckObject::StartCheckObjectTask) Running a check '"+c.GetName()+"'")

  exit := 0
  var output string
  var messageError string
  //Exit codes
  // OK: 0
  // WARN: 1
  // ERROR: 2
  // UNKNOWN: other (-1)
  cmdSplitted := strings.SplitN(c.GetCommand()," ",2)
  time_init := time.Now()
  out, err := exec.Command(cmdSplitted[0],strings.Split(cmdSplitted[1]," ")...).Output()
  elapsedtime := time.Since(time_init)

  // When the exec has exit code, these code is achived. If is not possible to achive it, then is set to '-1', the unknown code.
  if err != nil {
    if exiterr, ok := err.(*exec.ExitError); ok {
      if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
        exit = status.ExitStatus()
        env.Output.WriteChDebug("(CheckObject::StartCheckObjectTask) Exit status for '"+c.GetName()+"': "+strconv.Itoa(exit))
        if exit > 2 || exit < 0 {
          exit = -1
        }
      } else {
        exit = -1
      }
    } else {
      messageError = "(CheckObject::StartCheckObjectTask) The task for '"+c.GetName()+"has ended with errors"
      exit = -1
    }
  }

  if len(out) > 0 { output = string(out[:len(out)-1])}
  _,sample := c.GenerateCheckSample(exit,output,elapsedtime, time.Duration(c.GetExpirationTime())*time.Second,c.GetTimestamp())

  if messageError != "" {
    env.Output.WriteChWarn(messageError)
    return errors.New(messageError), sample
  } else {
    return nil, sample
  }
}
//
// GenerateCheckSample: method prepares the system to gather check's data
func (c *CheckObject) GenerateCheckSample(e int, o string, elapsedtime time.Duration, expirationtime time.Duration, timestamp int64) (error, *sample.CheckSample) {
  env.Output.WriteChDebug("(CheckObject::GenerateCheckSample) CheckSample for '"+c.GetName()+"'")
  checkEngine := env.GetCheckEngine().(*CheckEngine)
  cs := new(sample.CheckSample)

  cs.SetCheck(c.GetName())
  cs.SetExit(e)
  cs.SetOutput(o)
  cs.SetElapsedTime(elapsedtime)
  cs.SetSampletime(time.Now())
  cs.SetExpirationTime(expirationtime)
  cs.SetTimestamp(timestamp)

  env.Output.WriteChDebug("(CheckObject::GenerateCheckSample) "+cs.String())

  // send the sample to CheckEngines's sendSample method to write its value into output channels
  checkEngine.sendSample(cs)

  return nil,cs
}

//
// Common methods
//----------------------------------------------------------------------------------------

// String: converts a CheckObject object to string
func (c *CheckObject) String() string {
  if err, str := utils.ObjectToJsonString(c); err != nil{
    return err.Error()
  } else{
    return str
  }
}

//######################################################################################################
