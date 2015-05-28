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
  "verdmell/utils"
)
//#
//#
//# Checkgroups struct:
//# struct for checkgroups definition
type Checkgroups struct{
  Checkgroup map[string] []string `json:"checkgroups"`
}

//#
//# Getters/Setters methods for Checks object
//#---------------------------------------------------------------------

//# SetName method sets the Name value for the Checkgroups object
func (c *Checkgroups) SetCheckgroup( cg map[string] []string) {
  env.Output.WriteChDebug("(Checkgroups::SetCheckgroup) Set checkgroup")
  c.Checkgroup = cg
}

//# GetName method returns the Name value for the Checkgroups object
func (c *Checkgroups) GetCheckgroup() map[string] []string {
    env.Output.WriteChDebug("(Checkgroups::GetCheckgroup) Get checkgroup")
    return c.Checkgroup
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//# GetCheckgroupByName: returns a check object gived a name
func (c *Checkgroups) GetCheckgroupByName(checkgroupname string) ([]string, error) {
  env.Output.WriteChDebug("(Checkgroups::GetCheckObject) Looking for the checkgroup '"+checkgroupname+"'")

  var err bool
  var checkgroup []string
  checkgroups := c.GetCheckgroup()

  if checkgroup, err = checkgroups[checkgroupname]; err == false {
    return nil, errors.New("(Checkgroups::GetCheckgroupByName) The check group '"+checkgroupname+"' has never been load before.")
  }

  return checkgroup, nil
}

//
//# ValidateCheckgroups: validate that all checks defined for the checkgroups exist
func (c *Checkgroups) ValidateCheckgroups(i interface{}) error {
  env.Output.WriteChDebug("(Checkgroups::ValidateCheckgroups) Validate checkgroups")
	errorChan := make(chan error)
	check := i.(*Checks)
	checks := check.GetCheck()

	validation := func(cs []string) {
		for _, checkName := range cs {
			if _,ok := checks[checkName]; !ok {
				errorChan <- errors.New("(Checkgroups::ValidateCheckgroups) The check '"+checkName+"' does not exist")
			}
		}
		errorChan <- nil
	}

	for _, checkNames := range c.GetCheckgroup() {
		go validation(checkNames)
	}

	for i := 0; i < len(c.GetCheckgroup()); i++ {
		select{
    		case err := <- errorChan:
    			if err != nil {
    				return err
    			}	
		}
  }
  close(errorChan)
	return nil
}

//
//# UnmarshalCheckgroups: load checkgroup from file to an checkgroups object
func UnmarshalCheckgroups(file string) *Checkgroups{
  env.Output.WriteChDebug("checkgroups (Checkgroups::UnmarshalCheckgroup)")

  c := new(Checkgroups)
  utils.LoadJSONFile(file, c)

  return c
}

//
//# RetrieveCheckgroups: method load checks from the setup to system
func RetrieveCheckgroups(folder string) *Checkgroups{

  groups := new(Checkgroups)
  files := utils.GetFolderFiles(folder)
  // sync channels
  groupsFromFileChan := make(chan map[string][]string)
  groupsChan := make(chan *map[string][]string)

  // goroutine for extract each check object from file
  retrieveCheckgroupsFromFile :=  func(f os.FileInfo) {
    checkgroups := make(map[string] []string)
    checkFile := folder+string(os.PathSeparator)+f.Name()
    env.Output.WriteChDebug("(Checkgroups::RetrieveCheckGroups) File found: "+checkFile)

    c := UnmarshalCheckgroups(checkFile)
    for checkgroupName, checks := range c.GetCheckgroup(){
      if _,exist := checkgroups[checkgroupName]; !exist{
        env.Output.WriteChInfo("(Checkgroups::RetrieveCheckGroups) Checkgroup '"+checkgroupName+"' defined")
        checkgroups[checkgroupName] = checks
      } else {
        env.Output.WriteChWarn("(Checkgroups::RetrieveCheckgroups) The Checkgroup '"+checkgroupName+"' has already defined")
      }
    }
    groupsFromFileChan <- checkgroups
  }
  // call the goroutine for each file
  for _, f := range files {
    go retrieveCheckgroupsFromFile(f)
  }
  // waiting for all groupsFileEndChan that will indicate that all files has been analized
  go func() {
    allCheckgroups := make([]map[string][]string,0)
    for i := len(files); i > 0; i--{
      checkgroups := <-groupsFromFileChan
      allCheckgroups = append(allCheckgroups,checkgroups)
    }
    close(groupsFromFileChan)
    checkgroups := mergeCheckgroups(allCheckgroups)
    groupsChan <- checkgroups
  }()

  checkgroups := <-groupsChan
  close(groupsChan)
  //set the content of checkgroups to be set to groups
  groups.SetCheckgroup(*checkgroups)
  return groups
}

//
//# mergeCheckgroups: method merges a set of maps to one map
func mergeCheckgroups(allCheckgroups []map[string][]string) *map[string][]string{
    checkgroups := make(map[string][]string)

    for i := 0; i < len(allCheckgroups); i++ {
      for checkgroupName, checks := range allCheckgroups[i] {
        if _,exist := checkgroups[checkgroupName]; !exist{
            env.Output.WriteChDebug("(Checkgroups::mergeCheckgroups) New checkgroup: "+checkgroupName)
            checkgroups[checkgroupName] = checks
        } else {
          env.Output.WriteChWarn("(Checkgroups::mergeCheckgroups) The Checkgroup '"+checkgroupName+"' has already defined")
        }
      }
    }
    return &checkgroups
}

//#
//# Common methods
//#---------------------------------------------------------------------

//
//# String: method converts a Checks object to string
func (c *Checkgroups) String() string {
  return utils.ObjectToJsonString(c)
}
//#######################################################################################################