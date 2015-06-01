/*
Service System management

The package 'service' is used by verdmell to manage the definied services. 

=> Is known as a service a set of checks. By default the same node is a service compound by all defined checks.

-ServiceSystem
-Services
-ServiceObject
*/
package service

import (
	"errors"
  "verdmell/utils"
)

//#
//#
//# Service struct:
//# Service defines a map to store the maps
type ServiceObject struct{
	Name string `json: "name"`
  Description string `json:"description"`
  Checks []string `json:"checks"`
}

func NewServiceObject(name string, desc string, checks []string) (error,*ServiceObject) {
	//define new ServiceObject instance
	serviceObj := new(ServiceObject)
	serviceObj.SetName(name)
	serviceObj.SetDescription(desc)
	serviceObj.SetChecks(checks)


	if err := serviceObj.ValidateServiceObject(); err != nil {
		env.Output.WriteChDebug("(ServiceObject::NewServiceObject) Service not properly defined '"+serviceObj.GetName()+"'")
		return err, nil
	}

	env.Output.WriteChDebug("(ServiceObject::NewServiceObject) "+serviceObj.String())
	return nil, serviceObj
}

//#
//# Getters/Setters methods for Service object
//#---------------------------------------------------------------------

//
//# SetName: method sets the Name value for the ServiceObject
func (s *ServiceObject) SetName(n string){
  	s.Name = n
}
//
//# SetChecks: method sets the Checks value for the ServiceObject
func (s *ServiceObject) SetChecks(c []string){
  	s.Checks = c
}
//
//# SetDescription: method sets the description value for the ServiceObject
func (s *ServiceObject) SetDescription(d string){
  s.Description = d
}

//
//# GetName: method return the Name value for the ServiceObject
func (s *ServiceObject) GetName() string{
    return s.Name
}
//
//# GetCheck: method return the Checks value for the ServiceObject
func (s *ServiceObject) GetChecks() []string{
    return s.Checks
}
//
//# GetDescription: method return the description value for the ServiceObject
func (s *ServiceObject) GetDescription() string {
  return s.Description
}

//#
//# Specific methods
//#---------------------------------------------------------------------

func (s *ServiceObject) ValidateServiceObject() error {
	checkMap := make(map[string] interface{})

	if len(s.GetChecks()) < 1 {
		err := errors.New("(ServiceObject::ValidateServiceObject) Service '"+s.GetName()+"' must have a defined check")
    return err
	}

	// transform the []string to a map to optimize the search during validation
	for _,check := range env.GetChecks(){
		checkMap[check] = nil
	}
	// validate that all checks were defined as check
	for _,check := range s.GetChecks(){
		if _,exist := checkMap[check]; !exist{
			err := errors.New("(ServiceObject::ValidateServiceObject) Service '"+s.GetName()+"' requires the undefined check '"+check+"'")
    	return err
		}
	}

	return nil
}

//
//# String: method converts a ServiceObject to string
func (s *ServiceObject) String() string {
  return utils.ObjectToJsonString(s)
}
//#######################################################################################################