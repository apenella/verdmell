/*
Check Engine management

The package 'check' is used by verdmell to manage the monitoring checks defined by user

-Checks
-CheckObject
-Checkgroups

*/
package check

import (
  "os"
  "errors"
  "strconv"
  "time"
  "verdmell/sample"
  "verdmell/utils"
)
//#
//#
//# Checks struct:
//# Checks is an struct where the checks are stored
type Checks struct{
  Check map[string] *CheckObject `json:"checks"`
}

//
// Getters/Setters methods for Checks object

//# SetCheck: methods sets the Check value for the Check object
func (c *Checks) SetCheck( ck map[string]*CheckObject) {
  c.Check = ck
}
//# GetCheck: methods gets the Check's value for a gived Check object
func (c *Checks) GetCheck() map[string]*CheckObject{
    return c.Check
}

//
// Specific methods

//
//# AddCheck: method add a new check to the Checks struct
func (c *Checks) AddCheck(obj *CheckObject) error {
  if _,exist := c.Check[obj.GetName()]; !exist{
    // env.Output.WriteChDebug("(Checks::AddCheck) New Check '"+obj.GetName()+"'")
    log.Debug("(Checks::AddCheck) New Check '"+obj.GetName()+"'")
    c.Check[obj.GetName()] = obj
  }
  return nil
}

//
//# ListCheckNames: returns an array with the check namess defined on Checks object
func (c *Checks) ListCheckNames() []string {
  var names []string
  for checkname, _ := range c.Check {
    // env.Output.WriteChDebug("(Checks::ListCheckNames) check name: "+checkname)
    log.Debug("(Checks::ListCheckNames) check name: "+checkname)
    // append each check name to names array
    names = append(names, checkname)
  }
  return names
}
//
//# IsDefined: return if a check is defined
func (c *Checks) IsDefined(name string) bool {
  _,exist := c.Check[name]
  return exist
}
//
//# GetCheckObjectByName: returns a check object gived a name
func (c *Checks) GetCheckObjectByName(checkname string) (error,*CheckObject) {
  var err bool
  checkObj := new(CheckObject)
  check := c.GetCheck()

  // env.Output.WriteChDebug("(Checks::GetCheckObjectByName) Looking for check '"+checkname+"'")
  log.Debug("(Checks::GetCheckObjectByName) Looking for check '"+checkname+"'")

  if checkObj, err = check[checkname]; err == false {
    return errors.New("(Checks::GetCheckObjectByName) The check '"+checkname+"' has never been load before."),nil
  }

  return nil,checkObj
}
//
//# ValidateChecks: ensures that all the CheckObject from the Checks object have been defined correctly.
func (c *Checks) ValidateChecks(i interface{}) error {
  errorChan := make(chan error)
  statusChan := make(chan bool)

  // validation is a goroutine that will validate one CheckObjet and will write the status into a channel
  validation := func(c *CheckObject) {
      if err := c.ValidateCheckObject(); err != nil {
        errorChan <- err
      } else {
        statusChan <- true
      }
  }

  // for each CheckObject is launched a validation function
  for _, checkObj := range c.GetCheck(){
    go validation(checkObj)
  }

  // the method waits for all the status. If an error occurs, the function returns it
  for i := 0; i < len(c.GetCheck()); i++ {
    select{
      case err := <-errorChan:
        close(errorChan)
        return err
      case <- statusChan:
        break
    }
  }

  close(statusChan)
  // if no error has been found, all CheckObjects have been defined correctly
  return nil
}
//
//# StartCheckTaskPools: start a pool for each check. For each pool are generated the check execution tasks
func (c *Checks) StartCheckTaskPools() error {
  // env.Output.WriteChDebug("(Checks::StartCheckTaskPools) Ready to start all pools for checks")
  log.Debug("(Checks::StartCheckTaskPools) Ready to start all pools for checks")

  sampleChan := make(chan *sample.CheckSample)
  defer close(sampleChan)
  statusChan := make(chan int)
  defer close(statusChan)
  doneChan := make(chan bool)
  defer close(doneChan)
  errChan := make(chan error)
  defer close(errChan)

  // go through all checks from Checks (map[string]CheckObject)
  for _,check := range c.GetCheck(){
    // runGraphList let to trace which objects are waiting to run
    runGraphList := make(map[string]interface{},0)
    // adding the current object into run list
    runGraphList[check.GetName()] = nil
    // each check will run under its own goroutine
    go func (obj *CheckObject, rgl map[string]interface{}) {
      // env.Output.WriteChDebug("(Checks::StartCheckTaskPools) Initializing tasks for '"+obj.GetName()+"''s pool")
      log.Debug("(Checks::StartCheckTaskPools) Initializing tasks for '"+obj.GetName()+"''s pool")
      if err, checksample := c.InitCheckTasks(obj, rgl); err == nil {
        sampleChan <- checksample
      } else {
        errChan <- err
      }
    }(check, runGraphList)
  }

  // waiting the CheckObjects results
  go func(){
    // env.Output.WriteChDebug("(Checks::StartCheckTaskPools) waiting for tasks "+strconv.Itoa(len(c.GetCheck()))+" to be finished")
    log.Debug("(Checks::StartCheckTaskPools) waiting for tasks "+strconv.Itoa(len(c.GetCheck()))+" to be finished")
    for i:= 0; i<len(c.GetCheck()); i++{
      select{
      case checksample := <-sampleChan:
        // env.Output.WriteChDebug("(Checks::StartCheckTaskPools)["+strconv.Itoa(int(checksample.GetTimestamp()))+"] End of task has been notified for '"+checksample.GetCheck()+"'")
        log.Debug("(Checks::StartCheckTaskPools)["+strconv.Itoa(int(checksample.GetTimestamp()))+"] End of task has been notified for '"+checksample.GetCheck()+"'")
        //
        // The samples will be send from the command invocation
        //
        // checkEngine := env.GetCheckEngine().(*CheckEngine)
        // checkEngine.sendSample(checksample)
      case err := <-errChan:
        // env.Output.WriteChDebug(err)
        log.Error(err)
      }
    }
    doneChan <- true
  }()

  // All checks had send its sample
  <-doneChan

  return nil
}
//
//# InitCheckTasks: is going to initialize a task for each check and its dependencies.
//# The task enqueu the check to be executed. All its dependencies have to be executed before it to be enqueued
func (c *Checks) InitCheckTasks(checkObj *CheckObject, runGraphList map[string]interface{}) (error, *sample.CheckSample) {
  // env.Output.WriteChDebug("(Checks::InitCheckTasks) Initializing the tasks for check '"+checkObj.GetName()+"'")
  log.Debug("(Checks::InitCheckTasks) Initializing the tasks for check '"+checkObj.GetName()+"'")
  var err error
  exitStatus := -1
  checkengine := env.GetCheckEngine().(*CheckEngine)
  checksample := new(sample.CheckSample)
  sampleChan := make(chan *sample.CheckSample)
  defer close(sampleChan)
  // statusChan := make(chan int)
  // defer close(statusChan)
  doneChan := make(chan int)
  defer close(doneChan)
  errChan := make(chan error)
  defer close(errChan)
  jumpDueErrChan := make(chan error)
  defer close(jumpDueErrChan)

  if len(checkObj.GetDepend()) > 0 {
    //
    // recursive condition: A dependency is found
    //
    // env.Output.WriteChDebug("(Checks::InitCheckTasks) The check '"+checkObj.GetName()+"' has dependencies")
    log.Debug("(Checks::InitCheckTasks) The check '"+checkObj.GetName()+"' has dependencies")
    // Add to the runGraphList all CheckObjects to be run before the current object
    // If the object has already exist into the runGraphList then exist a cycle dependency.
    // The current object couldn't exist to it's dependency graph
    for _,d := range checkObj.GetDepend(){
      go func(dep string, rgl map[string]interface{}){
        // validate that the check doesn't already exist into list
        if _,exist := rgl[dep]; exist {
          // if it exist an error is launch for this execution branch
          //env.Output.WriteChError(append([]interface{}{"(Checks::InitCheckTasks) ",dep,checkObj.GetName()},rgl))
          log.Error(append([]interface{}{"(Checks::InitCheckTasks) ",dep,checkObj.GetName()},rgl))
          jumpDueErrChan <- errors.New("(Checks::InitCheckTasks) Your defined check has a cycle dependency for '"+dep+"'. Detected while running '"+checkObj.GetName()+"'.")
        } else {
          // get a CheckObject by its name
          if err,obj := checkengine.GetCheckObjectByName(dep); err == nil {
            // env.Output.WriteChDebug("(Checks::InitCheckTasks) The check '"+checkObj.GetName()+"' depends to '"+dep+"'")
            log.Debug("(Checks::InitCheckTasks) The check '"+checkObj.GetName()+"' depends to '"+dep+"'")
            // the current check must be marked into runGraphList
            rgl[d] = nil
            if err, sampleDedend := c.InitCheckTasks(obj, rgl); err != nil {
              errChan <- err
            } else {
              sampleChan <- sampleDedend
            }
          } else {
            // return the error in case the GetCHeckObjectByName returns an error
            // if an undefined CheckObject is defined such a dependency one, jump it
            jumpDueErrChan <- err
          }
        }
      }(d,runGraphList)
    }

    // gather the results for the depended check
    go func(){
      exitStatus := -1
      for i:=0; i < len(checkObj.GetDepend());i++{
        select{
          case err = <- errChan:
            // env.Output.WriteChError(err)
            log.Error(err)
            exitStatus = 4
          case err = <-jumpDueErrChan:
            // env.Output.WriteChError(err)
            log.Error(err)
            exitStatus = 4
          case s := <-sampleChan:
            //Exit codes
            // OK: 0
            // WARN: 1
            // ERROR: 2
            // UNKNOWN: others (-1)
            //
            // exitStatus calculates the task status throughout dependency task execution
            if exitStatus < s.GetExit() {
              exitStatus = s.GetExit()
            }
        }
      }
      doneChan <- exitStatus
    }()

    // waiting the command execution
    exitStatus = <-doneChan

    // once all dependent checks have been executed the current object is executed
    if exitStatus != 2 {
      // env.Output.WriteChDebug("(Checks::InitCheckTasks) The '"+checkObj.GetName()+"''s dependencies has been already executed")
      log.Debug("(Checks::InitCheckTasks) The '"+checkObj.GetName()+"''s dependencies has been already executed")
      // delete the check to runGraphList
      delete(runGraphList,checkObj.GetName())
      // queue the object to be run
      checkObj.EnqueueCheckObject()
      // Once the task are queued and executed the result is sent using the CheckObject's SampleChan
      checksample = <-checkObj.SampleChan
      exitStatus = checksample.GetExit()
      // env.Output.WriteChDebug("(Checks::InitCheckTasks) Received a check status for '"+checkObj.GetName()+"': '"+strconv.Itoa(exitStatus)+"'")
      log.Debug("(Checks::InitCheckTasks) Received a check status for '"+checkObj.GetName()+"': '"+strconv.Itoa(exitStatus)+"'")
    }else{
      outputMessage := "Wrong status for '"+checkObj.GetName()+"' because it depends to another check with "+sample.Itoa(exitStatus)+" status"
      // env.Output.WriteChWarn("(Checks::InitCheckTasks) "+outputMessage)
      log.Warn("(Checks::InitCheckTasks) "+outputMessage)
      _,checksample = checkObj.GenerateCheckSample(-1,outputMessage,time.Duration(0)*time.Second, time.Duration(0)*time.Second, checkObj.GetTimestamp())

      go func() {
        // env.Output.WriteChWarn("(Checks::InitCheckTasks) Countdown for '"+checkObj.GetName()+"'")
        log.Warn("(Checks::InitCheckTasks) Countdown for '"+checkObj.GetName()+"'")
        timeout := time.After(time.Duration(checkObj.GetInterval()) * time.Second)
        for{
          select{
          case <-timeout:
            c.InitCheckTasks(checkObj,runGraphList)
          }
        }
      }()
    }
  }else{
    //
    // recursive ending condition: No dependency is found
    //
    // env.Output.WriteChDebug("(Checks::InitCheckTasks) The check '"+checkObj.GetName()+"' hasn't dependencies")
    log.Debug("(Checks::InitCheckTasks) The check '"+checkObj.GetName()+"' hasn't dependencies")
    // delete the check to runGraphList
    delete(runGraphList,checkObj.GetName())
    // queue the object to be run
    checkObj.EnqueueCheckObject()
    // Once the task are queued and executed the result is sent using the CheckObject's SampleChan
    checksample = <-checkObj.SampleChan
    exitStatus = checksample.GetExit()
    // env.Output.WriteChDebug("(Checks::InitCheckTasks) Received a check status for '"+checkObj.GetName()+"': '"+strconv.Itoa(exitStatus)+"'")
    log.Debug("(Checks::InitCheckTasks) Received a check status for '"+checkObj.GetName()+"': '"+strconv.Itoa(exitStatus)+"'")
  }

  return nil, checksample
}
//
//# UnmarshalCheck: get the json content from a file and field an Checks object on it.
//  The method requieres a file path.
//  The method returns a pointer to Checks object
func UnmarshalCheck(file string) *Checks {
  // env.Output.WriteChDebug("(Checks::UnmarshalCheck)")
  log.Debug("(Checks::UnmarshalCheck)")
  c := new(Checks)
  // extract the content from the file and dumps it on the CHecks object
  if err := utils.LoadJSONFile(file, c); err != nil {
    // env.Output.WriteChError("(Checks::UnmarshalCheck) The input file '"+file+"' has an invalid json structure")
    log.Error("(Checks::UnmarshalCheck) The input file '"+file+"' has an invalid json structure")
  }

  return c
}
//
//# RetrieveChecks: gets all the files found on checks folder and generate one Checks object with all this CheckObject defined.
func RetrieveChecks(folder string) *Checks{
  check := new(Checks)
  // checks will contain all the CheckObject definition
  checks := make(map[string]*CheckObject)
  // files is an array with all files found inside the folder
  files := utils.GetFolderFiles(folder)
  // sync channel
  checkObjChan := make(chan *CheckObject)
  checkFileEndChan := make(chan bool)
  allChecksGetChan := make(chan bool)
  done := make(chan *Checks)

  // goroutine for extract each check object from file
  retrieveChecksFromFile := func(f os.FileInfo) {
    checkFile := folder+string(os.PathSeparator)+f.Name()
    // env.Output.WriteChDebug("(Checks::RetrieveChecks) File found: "+checkFile)
    log.Debug("(Checks::RetrieveChecks) File found: "+checkFile)

    c := UnmarshalCheck(checkFile)

    if len(c.GetCheck()) == 0 { env.Output.WriteChWarn("(Checks::RetrieveChecks) You should review the file "+checkFile+", no check has been load from it") }
    for checkName, checkObj := range c.GetCheck(){
      queue := make(chan *CheckObject)
      sample := make(chan *sample.CheckSample)

      // the CheckObject Name may be set because in the json file comes as a key
      checkObj.SetName(checkName)
      // the CheckObject Queue may be set to proceed the execution requests
      checkObj.SetTaskQueue(queue)
      // the CheckObject StatusChan may be set to proceed the execution requests
      checkObj.SetSampleChan(sample)
      // the CheckObject Timestamp to 0
      checkObj.SetTimestamp(0)

      if checkObj.GetExpirationTime() < 0 {
        // env.Output.WriteChDebug("(Checks::RetrieveChecks) The expiration time for '"+checkObj.GetName()+"' has not been defined properly and will be overwritten")
        log.Warn("(Checks::RetrieveChecks) The expiration time for '"+checkObj.GetName()+"' has not been defined properly and will be overwritten")
        checkObj.SetExpirationTime(300)
      }

      if checkObj.GetInterval() < checkObj.GetExpirationTime() {
        // env.Output.WriteChDebug("(Checks::RetrieveChecks) The interval time for '"+checkObj.GetName()+"' has not been defined properly and will be overwritten")
        log.Warn("(Checks::RetrieveChecks) The interval time for '"+checkObj.GetName()+"' has not been defined properly and will be overwritten")
        checkObj.SetInterval(checkObj.GetExpirationTime())
      }

      // sending the CheckObject to be stored
      checkObjChan <- checkObj
      // env.Output.WriteChDebug("(Checks::RetrieveChecks) Check '"+checkName+"' defined")
      // env.Output.WriteChDebug("(Checks::RetrieveChecks) '"+checkObj.String()+"'")
    }
    // a message is send when all CheckObject defined into a file have been sent to store
    checkFileEndChan <- true
  }
  // call the goroutine for each file
  for _, f := range files {
    go retrieveChecksFromFile(f)
  }
  // waiting for all checkFileEndChan that will indicate that all files has been analized
  go func() {
    for i := len(files); i > 0; i--{
      <-checkFileEndChan
    }
    defer close(checkFileEndChan)
    allChecksGetChan <- true
  }()
  // store all CheckObjects sent. Once the allChecksGetChan channel gets a message the goroutine will assume that all CheckOjects has been sent
  go func() {
    check := new(Checks)
    allChecksGet := false
    for ;!allChecksGet;{
      select{
        // get a CheckObject object
        case obj := <- checkObjChan:
          // env.Output.WriteChDebug("(Checks::RetrieveChecks::routine) New check to register '"+obj.GetName()+"'")
          log.Debug("(Checks::RetrieveChecks::routine) New check to register '"+obj.GetName()+"'")
          if _,exist := checks[obj.GetName()]; !exist{
            checks[obj.GetName()] = obj
          }
        // ending message
        case allChecksGet = <-allChecksGetChan:
          check.SetCheck(checks)
          done <-check
          defer close(checkObjChan)
          defer close(allChecksGetChan)
      }
    }
  }()
  // the main routine will wait for the work to be done
  check = <-done
  defer close(done)

  return check
}
//
//# Itoa: method transform a integer to status string
func Itoa(i int) string {
  switch(i){
  case 0:
    return "OK"
  case 1:
    return "WARN"
  case 2:
    return "ERROR"
  default:
    return "UNKNOWN"
  }
}

//#
//# Common methods
//#---------------------------------------------------------------------

//# String: convert a Checks object to string
func (c *Checks) String() string {
if err, str := utils.ObjectToJsonString(c); err != nil{
    return err.Error()
  } else{
    return str
  }
}
//#######################################################################################################
