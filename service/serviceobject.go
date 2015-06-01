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
	"verdmell/sample"
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
  Status int `json:"status"`
  inputSampleChan chan *sample.CheckSample `json:"-"`
}

func NewServiceObject(name string, desc string, checks []string) (error,*ServiceObject) {
	//define new ServiceObject instance
	serviceObj := new(ServiceObject)
	serviceObj.SetName(name)
	serviceObj.SetDescription(desc)
	serviceObj.SetChecks(checks)
	serviceObj.SetStatus(-1)

	if err := serviceObj.ValidateServiceObject(); err != nil {
		env.Output.WriteChDebug("(ServiceObject::NewServiceObject) Service not properly defined '"+serviceObj.GetName()+"'")
		return err, nil
	}

	go func() {
		serviceObj.StartServiceObjectSampleChannel()
	}()
	
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
//# SetStatus: method sets the Status value for the ServiceObject
func (s *ServiceObject) SetStatus(status int){
  s.Status = status
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
//
//# SetStatus: method sets the Status value for the ServiceObject
func (s *ServiceObject) GetStatus() int{
  return s.Status
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//#ValidateServiceObject: validates SericeObject
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
//# StartServiceObjectCheckSampleInput: methot the serviceObject to receive samples and calculates the service status
func (s *ServiceObject) StartServiceObjectSampleChannel(){
	sampleChan := make(chan *sample.CheckSample)
	defer close(sampleChan)
	s.inputSampleChan = sampleChan
	
	env.Output.WriteChDebug("(ServiceObject::StartServiceObjectCheckSampleInput) Waiting samles for service '"+s.GetName()+"'")
	for {
		select{
		case sam := <- s.inputSampleChan:
			env.Output.WriteChDebug("(ServiceObject::StartServiceObjectCheckSampleInput) New sample arrived for '"+sam.GetCheck()+"' to service '"+s.GetName()+"'")
			s.CalculateStatusForService(sam.GetExit())
		}
	}
}

//
//# SendToSampleChannel: method sends a sample to the sample channel
func (s *ServiceObject) SendToSampleChannel(sample *sample.CheckSample){
	env.Output.WriteChDebug("(ServiceObject::SendToSampleChannel) New sample to send for '"+sample.GetCheck()+"' to service '"+s.GetName()+"'")
	s.inputSampleChan <- sample
}

//
//# SendToSampleChannel: method sends a sample to the sample channel
func (s *ServiceObject) CalculateStatusForService(newStatus int){
	//Exit codes
  // OK: 0
  // WARN: 1
  // ERROR: 2
  // UNKNOWN: others (-1)
  //
	currentStatus := s.GetStatus()
	//exitStatus calculates the task status throughout dependency task execution
  if currentStatus < newStatus {
  	env.Output.WriteChDebug("(ServiceObject::CalculateStatusForService) Service '"+s.GetName()+"' has changed its status to '"+sample.Itoa(newStatus)+"'")
		s.SetStatus(newStatus)
  }
}


//#
//# Common methods
//#---------------------------------------------------------------------

//
//# String: method converts a ServiceObject to string
func (s *ServiceObject) String() string {
  return utils.ObjectToJsonString(s)
}
//#######################################################################################################