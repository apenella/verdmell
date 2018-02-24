/*
Service System management

The package 'service' is used by verdmell to manage the definied services. 

=> Is known as a service a set of checks. By default the same node is a service compound by all defined checks.

-ServiceEngine
-Services
-ServiceObject
-ServiceResult
*/
package service

import (
	"time"
	"errors"
	"strconv"
	"verdmell/check"
	"verdmell/sample"
	"verdmell/utils"
)

//#
//#
//# Service struct:
//# Service defines a map to store the maps
type ServiceObject struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Checks []string `json:"checks"`
	Status int `json:"status"`
	Timestamp int64 `json:"timestamp"`
	inputSampleChan chan *sample.CheckSample `json:"-"`
	checksStatusCache map[string] int	`json:"-"`
}

func NewServiceObject(name string, desc string, checks []string) (error,*ServiceObject) {
	//define new ServiceObject instance
	serviceObj := new(ServiceObject)
	serviceObj.SetName(name)
	serviceObj.SetDescription(desc)
	serviceObj.SetChecks(checks)
	serviceObj.SetStatus(-1)
	serviceObj.SetTimestamp(0)

	if err := serviceObj.ValidateServiceObject(); err != nil {
		env.Output.WriteChDebug("(ServiceObject::NewServiceObject) Service not properly defined '"+serviceObj.GetName()+"'")
		return err, nil
	}

	go func() {
		serviceObj.StartReceiver()
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
//# SetTimestamp: attribute from CheckObject
func (s *ServiceObject) SetTimestamp(t int64) {
  env.Output.WriteChDebug("(ServiceObject::SetTimestamp)")
  s.Timestamp = t
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
//
//# GetTimestamp: attribute from CheckObject
func (s *ServiceObject) GetTimestamp() int64 {
  env.Output.WriteChDebug("(ServiceObject::GetTimestamp)")
  return s.Timestamp
}


//#
//# Specific methods
//#---------------------------------------------------------------------

//
//#CopyServiceObject: method copies the service ojbect
func (s *ServiceObject) CopyServiceObject() *ServiceObject {
	if s == nil {
		return nil
	}

	serviceObj := new(ServiceObject)
	serviceObj.SetName(s.GetName())
	serviceObj.SetDescription(s.GetDescription())
	serviceObj.SetChecks(s.GetChecks())
	serviceObj.SetStatus(s.GetStatus())
	serviceObj.SetTimestamp(s.GetTimestamp())
	return serviceObj
}

//
//#WaitAllSamples: method waits that all checks send at least one sample
func (s *ServiceObject) WaitAllSamples(seconds int) *ServiceObject {
	env.Output.WriteChDebug("(ServiceObject::WaitAllSamples) Waiting "+strconv.Itoa(seconds)+"s")

	arrived := make(chan bool)
	defer close(arrived)

	go func() {
		allArrived := (len(s.GetChecks()) == len(s.checksStatusCache))
		if allArrived { arrived <- true }

		for;!allArrived;{
			env.Output.WriteChDebug("(ServiceObject::WaitAllSamples) "+strconv.Itoa(len(s.checksStatusCache))+ " samples arrived from "+ strconv.Itoa(len(s.GetChecks())) )
			if len(s.GetChecks()) == len(s.checksStatusCache) {
				env.Output.WriteChDebug("(ServiceObject::WaitAllSamples) all arrived")
				allArrived = (len(s.GetChecks()) == len(s.checksStatusCache))
				arrived <- allArrived
			}
		}
	}()

	// There is a timeout to avoid long wait...
	timeout := time.After(time.Duration(seconds) * time.Second)
	select{
	case <-arrived:
		env.Output.WriteChDebug("(ServiceObject::WaitAllSamples) All sample arrived...")
	case <-timeout:
		env.Output.WriteChWarn("(ServiceObject::WaitAllSamples) Samples for '"+s.GetName()+"' has not already arrived after "+strconv.Itoa(seconds)+" seconds")
	}

	return s
}

//
//#ValidateServiceObject: validates SericeObject
func (s *ServiceObject) ValidateServiceObject() error {
	checkMap := make(map[string] interface{})

	if len(s.GetChecks()) < 1 {
		err := errors.New("(ServiceObject::ValidateServiceObject) Service '"+s.GetName()+"' must have a defined check")
		return err
	}

	// transform the []string to a map to optimize the search during validation
	checkengine := env.GetCheckEngine().(*check.CheckEngine)
	for _,check := range checkengine.ListCheckNames() {
		checkMap[check] = nil
	}
	// validate that all checks were defined as check
	for _,check := range s.GetChecks(){
		if _,exist := checkMap[check]; !exist{
			err := errors.New("(ServiceObject::ValidateServiceObject) Service '"+s.GetName()+"' requires the undefined check '"+check+"'")
			return err
		}
	}
	env.Output.WriteChDebug("(ServiceObject::ValidateServiceObject) The service '"+s.GetName()+"' has been properly validated")
	return nil
}
 
//
//# StartReceiver: methot the serviceObject to receive samples and calculates the service status
func (s *ServiceObject) StartReceiver(){
	sampleChan := make(chan *sample.CheckSample)
	defer close(sampleChan)
	checksStatusCache := make(map[string] int)

	s.inputSampleChan = sampleChan
	s.checksStatusCache = checksStatusCache
	
	env.Output.WriteChDebug("(ServiceObject::StartReceiver) Waiting samles for service '"+s.GetName()+"'")
	for {
		select{
		case sam := <- s.inputSampleChan:
			env.Output.WriteChDebug("(ServiceObject::StartReceiver) New sample arrived for '"+sam.GetCheck()+"' to service '"+s.GetName()+"'")

			statusCachedValue, exist := s.checksStatusCache[sam.GetCheck()]

			if !exist || statusCachedValue != sam.GetExit() {
				env.Output.WriteChDebug("(ServiceObject::StartReceiver) The '"+sam.GetCheck()+"' status has changed, and service status have to be calculated.")
				s.CalculateStatusForService(sam)
				// start a new routine that will send the service object to interested ones
				go func() {
					serviceEngine := env.GetServiceEngine().(*ServiceEngine)
					serviceEngine.SendData(s)					
				}()
			}

		}
	}
}

//
//# SendToSampleChannel: method sends a sample to the sample channel
func (s *ServiceObject) RecevieData(sample *sample.CheckSample){
	env.Output.WriteChDebug("(ServiceObject::ReceiveData) New sample to send for '"+sample.GetCheck()+"' to service '"+s.GetName()+"'")
	s.inputSampleChan <- sample
}
//
//# SendToSampleChannel: method sends a sample to the sample channel
func (s *ServiceObject) CalculateStatusForService(sam *sample.CheckSample){
	env.Output.WriteChDebug("(ServiceObject::CalculateStatusForService) Calculating '"+sam.GetCheck()+"' status to service '"+s.GetName()+"'")
	//Exit codes
	// OK: 0
	// WARN: 1
	// ERROR: 2
	// UNKNOWN: others (-1)
	//
	//currentStatus := s.GetStatus()
	currentStatus := -1

	s.checksStatusCache[sam.GetCheck()] = sam.GetExit()

	for _,status := range s.checksStatusCache {
		//exitStatus calculates the task status throughout dependency task execution
		if currentStatus < status {
			currentStatus = status
			env.Output.WriteChDebug("(ServiceObject::CalculateStatusForService) Service '"+s.GetName()+"' has changed its status to '"+sample.Itoa(status)+"'")
		}
	}
	s.SetTimestamp(s.GetTimestamp()+1)
	s.SetStatus(currentStatus)
}

//#
//# Common methods
//#---------------------------------------------------------------------

//
//# String: method converts a ServiceObject to string
func (s *ServiceObject) String() string {
	if err, str := utils.ObjectToJsonString(s); err != nil{
		return err.Error()
	} else{
		return str
	}
}
//#######################################################################################################