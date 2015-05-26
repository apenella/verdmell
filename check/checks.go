/*
Check system management

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
  Check map[string] CheckObject `json:"checks"`
}
//#
//# Getters/Setters methods for Checks object
//#---------------------------------------------------------------------

//# SetCheck: methods sets the Check value for the Check object
func (c *Checks) SetCheck( ck map[string]CheckObject) {
  c.Check = ck
}
//# GetCheck: methods gets the Check's value for a gived Check object
func (c *Checks) GetCheck() map[string]CheckObject{
    return c.Check
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# AddCheck: method add a new check to the Checks struct
func (c *Checks) AddCheck(obj CheckObject) error { 
  if _,exist := c.Check[obj.GetName()]; !exist{
    env.Output.WriteChDebug("(Checks::AddCheck) New Check '"+obj.GetName()+"'")
    c.Check[obj.GetName()] = obj
  }
  return nil
}

//
//# GetCheckNames: returns an array with the check namess defined on Checks object 
func (c *Checks) GetCheckNames() []string {
  var names []string
  for checkname, _ := range c.Check {
    env.Output.WriteChDebug("(Checks::GetCheckNames) check name: "+checkname)
    // append each check name to names array
    names = append(names, checkname)
  }
  return names
}
//
//# GetCheckObject: returns a check object gived a name
func (c *Checks) GetCheckObjectByName(checkname string) (*CheckObject, error) {
  var err bool
  checkObj := new(CheckObject)
  check := c.GetCheck()

  env.Output.WriteChDebug("(Checks::GetCheckObject) Looking for the check '"+checkname+"'")

  if *checkObj, err = check[checkname]; err == false {
    return nil, errors.New("(Checks::GetCheckObject) The checkname '"+checkname+"' has never been load before.")
  }

  return checkObj, nil
}
//
//# ValidateChecks: ensures that all the CheckObject from the Checks object have been defined correctly. 
func (c *Checks) ValidateChecks(i interface{}) error {
  errorChan := make(chan error)
  statusChan := make(chan bool)

  // validation is a goroutine that will validate one CheckObjet and will write the status into a channel
  validation := func(c CheckObject) {
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
func (c *Checks) StartCheckTaskPools(ss *sample.SampleSystem) (error, int){
  env.Output.WriteChDebug("(Checks::StartCheckTaskPools) Ready to start all pools for checks")

  sampleChan := make(chan *sample.CheckSample)
  defer close(sampleChan)
  statusChan := make(chan int)
  defer close(statusChan)
  doneChan := make(chan int)
  defer close(doneChan)
  errChan := make(chan error)
  defer close(errChan)

  exitStatus := -1

  // go over Check attriutes from Checks (map[string]CheckObject)
  for _,check := range c.GetCheck(){
    // runGraphList let to trace which objects are waiting to run
    runGraphList := make(map[string]interface{},0)
    // adding the current object into run list
    runGraphList[check.GetName()] = nil
    // each check will run under its own goroutine
    go func (o CheckObject, rgl map[string]interface{}) {
      env.Output.WriteChDebug("(Checks::StartCheckTaskPools) Initializing tasks for '"+o.GetName()+"''s pool")
      if err, checksample := c.InitCheckTasks(o, rgl); err == nil {
        sampleChan <- checksample
      } else {
        errChan <- err
      }
    }(check, runGraphList)
  }

  // waiting the CheckObjects results
  go func(){
    exitStatus := -1

    for i:= 0; i<len(c.GetCheck()); i++{
      select{
      case checksample := <-sampleChan:
        env.Output.WriteChDebug("(Checks::StartCheckTaskPools) Check status received: '"+strconv.Itoa(checksample.GetExit())+"'")
        ss.AddSample(checksample)
        //Exit codes
        // OK: 0
        // WARN: 1
        // ERROR: 2
        // UNKNOWN: others (-1)
        //
        // exitStatus calculates the task status throughout dependency task execution
        if exitStatus < checksample.GetExit(){
          exitStatus =  checksample.GetExit()
        }
      case err := <-errChan:
        env.Output.WriteChDebug(err)
        exitStatus = 2
      }
    }
    doneChan <- exitStatus
  }()

  exitStatus = <-doneChan
  env.Output.WriteChDebug("(Checks::StartCheckTaskPools) Check task pool status: '"+strconv.Itoa(exitStatus)+"'")

  return nil, exitStatus
}
//
//# InitCheckTasks: is going to initialize a task for each check and its dependencies. 
//#  The task enqueu the check to be check but it has dependencies they have to be enqueued befor
func (c *Checks) InitCheckTasks(checkObj CheckObject, runGraphList map[string]interface{}) (error, *sample.CheckSample) {
  env.Output.WriteChDebug("(Checks::InitCheckTasks) Initializing the tasks for the check '"+checkObj.GetName()+"'")

  var err error
  exitStatus := -1

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
    env.Output.WriteChDebug("(Checks::InitCheckTasks) The check '"+checkObj.GetName()+"' has dependencies")
    // Add to the runGraphList all CheckObjects to be run before the current object
    // If the object has already exist into the runGraphList then exist a cycle dependency. 
    // The current object couldn't exist to it's dependency graph
    for _,d := range checkObj.GetDepend(){  
      // validate that the check doesn't already exist into list
      if _,exist := runGraphList[d]; exist {
        // if it exist an error is launch for this execution branch
        go func() {
          jumpDueErrChan <- errors.New("(Checks::InitCheckTasks) Your defined check has a cycle dependency for '"+d+"'. Detected while running '"+checkObj.GetName()+"'.")
        }()

      } else {
        // get a CheckObject by its name
        if co, err := c.GetCheckObjectByName(d); err == nil {
          env.Output.WriteChDebug("(Checks::InitCheckTasks) The check '"+checkObj.GetName()+"' has a dependency to '"+d+"'")
          // runGraph for each object an wait for its response
          go func (o CheckObject) {
            // the current check must be marked into runGraphList
            runGraphList[d] = nil 
            if err, sampleDedend := c.InitCheckTasks(o, runGraphList); err != nil {
              errChan <- err
            } else {
              sampleChan <- sampleDedend
            }
          }(*co)
        } else {
          // return the error in case the GetCHeckObjectByName returns an error
          // if an undefined CheckObject is defined such a dependency one, jump it
          go func() { jumpDueErrChan <- err }()
        }
      }
    }

    // gather the results for the depended check
    go func(){
      exitStatus := -1
      for i:=0; i < len(checkObj.GetDepend());i++{
        select{
          case err = <- errChan:
            env.Output.WriteChError(err)
            exitStatus = 4
          case err = <-jumpDueErrChan:
            env.Output.WriteChError(err)
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
      checkObj.EnqueueCheckObject()
      checksample = <-checkObj.SampleChan
      exitStatus = checksample.GetExit()
      env.Output.WriteChDebug("(Checks::runCheckGraph) Received a check status for '"+checkObj.GetName()+"': '"+strconv.Itoa(exitStatus)+"'")
    }else{
      outputMessage := "Wrong status for '"+checkObj.GetName()+"' due dependency issue. Dependency status: '"+strconv.Itoa(exitStatus)+"'"
      env.Output.WriteChDebug("(Checks::runCheckGraph) "+outputMessage)
      _,checksample = checkObj.GenerateCheckSample(-1,outputMessage,time.Duration(0)*time.Second, time.Duration(0)*time.Second)
      //exitStatus = checksample.GetExit()
    }
  }else{
    //
    // recursive ending condition: No dependency is found
    //
    env.Output.WriteChDebug("(Checks::runCheckGraph) The check '"+checkObj.GetName()+"' hasn't dependencies")
    // delete the check to runGraphList
    delete(runGraphList,checkObj.GetName())
    //run the check
    checkObj.EnqueueCheckObject()
    checksample = <-checkObj.SampleChan
    exitStatus = checksample.GetExit() 
    env.Output.WriteChDebug("(Checks::runCheckGraph) Received a check status for '"+checkObj.GetName()+"': '"+strconv.Itoa(exitStatus)+"'")
  }

  return nil, checksample
}
//
//# UnmarshalCheck: get the json content from a file and field an Checks object on it.
//  The method requieres a file path.
//  The method returns a pointer to Checks object
func UnmarshalCheck(file string) *Checks{
  env.Output.WriteChDebug("(Checks::UnmarshalCheck)")

  c := new(Checks)
  // extract the content from the file and dumps it on the CHecks object
  utils.LoadJSONFile(file, c)

  return c
}
//
//# RetrieveChecks: gets all the files found on checks folder and generate one Checks object with all this CheckObject defined.
func RetrieveChecks(folder string) *Checks{
  check := new(Checks)
  // checks will contain all the CheckObject definition
  checks := make(map[string]CheckObject)
  // files is an array with all files found inside the folder
  files := utils.GetFolderFiles(folder)
  // sync channel
  checkObjChan := make(chan CheckObject)
  checkFileEndChan := make(chan bool)
  allChecksGetChan := make(chan bool)
  done := make(chan *Checks)

  // goroutine for extract each check object from file
  retrieveChecksFromFile := func(f os.FileInfo) {
    checkFile := folder+string(os.PathSeparator)+f.Name()
    env.Output.WriteChDebug("(Checks::RetrieveChecks) File found: "+checkFile)
    
    c := UnmarshalCheck(checkFile)

    if len(c.GetCheck()) == 0 { env.Output.WriteChWarn("(Checks::RetrieveChecks) You should review the file "+checkFile+", no check has been load from it") }
    for checkName, checkObj := range c.GetCheck(){
      queue := make(chan *CheckObject)
      sample := make(chan *sample.CheckSample)

      if checkObj.GetExpirationTime() < 0 {
        checkObj.SetExpirationTime(300)
      }

      // the CheckObject Name may be set because in the json file comes as a key
      checkObj.SetName(checkName)
      // the CheckObject Queue may be set to proceed the execution requests
      checkObj.SetTaskQueue(queue)
      // the CheckObject StatusChan may be set to proceed the execution requests
      checkObj.SetSampleChan(sample)

      // sending the CheckObject to be stored
      checkObjChan <- checkObj
      env.Output.WriteChInfo("(Checks::RetrieveChecks) Check '"+checkName+"' defined")
      env.Output.WriteChDebug("(Checks::RetrieveChecks) '"+checkObj.String()+"'")
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
    var check Checks
    allChecksGet := false
    for ;!allChecksGet;{
      select{
        // get a CheckObject object
        case obj := <- checkObjChan:
          env.Output.WriteChDebug("(Checks::RetrieveChecks::routine) New check to store: "+obj.GetName())
          if _,exist := checks[obj.GetName()]; !exist{
            checks[obj.GetName()] = obj
          }
        // ending message
        case allChecksGet = <-allChecksGetChan:
          check.SetCheck(checks)
          done <-&check
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
  return utils.ObjectToJsonString(c)
}
//#######################################################################################################