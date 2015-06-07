/*
Service System management

The package 'service' is used by verdmell to manage the definied services. 

=> Is known as a service a set of checks. By default the same node is a service compound by all defined checks.

-ServiceSystem
-Services
-ServiceObject
-ServiceResult
*/
package service

import (
  "verdmell/environment"
  "verdmell/utils"
  "verdmell/sample"
)

//
var env *environment.Environment

//#
//#
//# ServiceSystem struct:
//# ServiceSystem defines a map to store the maps
type ServiceSystem struct{
  Ss *Services `json:"servicesroot"`
  inputSampleChan chan *sample.CheckSample `json:"-"`
}

//
//# NewCheckSystem: return a Checksystem instance to be run
func NewServiceSystem(e *environment.Environment) (error, *ServiceSystem){
  e.Output.WriteChDebug("(ServiceSystem::NewServiceSystem)")
  sys := new(ServiceSystem)
  
  // set the environment attribute
  env = e

  // folder that contains services definitions
  folder := env.Setup.Servicesfolder
  // Get defined services
  srv := RetrieveServices(folder)

  // set the attribute CheckServiceMapReduce
  srv.GenerateCheckServices()

  if err := srv.ValidateServices(); err != nil {
    return err, nil
  }
  sys.SetServices(srv)

  // Set description for default service
  desc := "Global services for node "+env.Setup.Hostname
  
  if err := sys.RegisterService(env.Setup.Hostname,desc,env.GetChecks()); err != nil{
    return err, nil
  }

  // start the sample receiver
  sys.StartServiceSystemReceiver()

  return nil, sys
}

//#
//# Getters/Setters methods for ServiceSystem object
//#---------------------------------------------------------------------

//
//# SetServices: methods sets the Services' value
func (s *ServiceSystem) SetServices(ss *Services) {
  s.Ss = ss
}
//
//# SetInputSampleChan: methods sets the inputSampleChan's value
func (s *ServiceSystem) SetInputSampleChan(c chan *sample.CheckSample) {
  s.inputSampleChan = c
}

//
//# GetServices: methods gets the Services' value
func (s *ServiceSystem) GetServices() *Services {
  return s.Ss
}
//
//# GetInputSampleChan: methods sets the inputSampleChan's value
func (s *ServiceSystem) GetInputSampleChan() chan *sample.CheckSample {
  return s.inputSampleChan
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# StartServiceSystem: method prepares the system to wait sample and calculate the results for services
func (s *ServiceSystem) StartServiceSystemReceiver() error {
  s.inputSampleChan = make(chan *sample.CheckSample)
  services := s.GetServices()

  env.Output.WriteChDebug("(ServiceSystem::StartServiceSystemReceiver) Starting sample receiver")
  go func() {
    defer close (s.inputSampleChan)
    for{
      select{
      case sample := <-s.inputSampleChan:
        env.Output.WriteChDebug("(ServiceSystem::StartServiceSystemReceiver) New sample received for '"+sample.GetCheck()+"'")
        _,servicesCheck := s.GetServicesForCheck(sample.GetCheck())
        for _,service := range servicesCheck {
          _,srv := services.GetServiceObject(service)
          env.Output.WriteChDebug("(ServiceSystem::StartServiceSystemReceiver) Sample for '"+sample.GetCheck()+"' belongs to '"+srv.GetName()+"'")
          go srv.SendToSampleChannel(sample)
        }
      }
    }
  }()
  return nil
}
//
//# SendSampleToServiceSystem: method prepares the system to wait sample and calculate the results for services
func (s *ServiceSystem) SendSampleToServiceSystem(sample *sample.CheckSample) {
  env.Output.WriteChDebug("(ServiceSystem::SendSampleToServiceSystem) Send sample "+sample.String())
  s.inputSampleChan <- sample
}
//
//# RegisterService: register a new service for ServiceSysem
func (s *ServiceSystem) RegisterService(name string, desc string, checks []string) error {
  env.Output.WriteChDebug("(ServiceSystem::RegisterService) New service to register '"+name+"'")
  var serviceObj *ServiceObject
  var err error

  if err, serviceObj = NewServiceObject(name, desc, checks); err != nil {
    return err
  }
  srv := s.GetServices()
  // add the service for node
  srv.AddServiceObject(serviceObj)
  // set the attribute CheckServiceMapReduce
  srv.GenerateCheckServices()

  return nil
}

//
//# GetAllServicesStatusHuman: converts a SampleSystem object to string
func (sys *ServiceSystem) GetAllServicesStatusHuman() (error, string) {
  env.Output.WriteChDebug("(ServiceSystem::GetAllServicesStatusHuman)")
  var str string
  var substr string
  var err error

  ss := sys.GetServices()

  for _,obj := range ss.GetServices(){
    if err, substr = sys.GetServicesStatusHuman(obj.GetName()); err != nil {
      return err, substr
    }
    str += substr
  }
  return nil, str
}
//
//# GetServicesStatusHuman: converts a SampleSystem object to string
func (sys *ServiceSystem) GetServicesStatusHuman(service string) (error ,string) {
  env.Output.WriteChDebug("(ServiceSystem::GetServicesStatusHuman) Get status for '"+service+"'")
  var obj *ServiceObject
  var err error

  ss := sys.GetServices()
  if err, obj = ss.GetServiceObject(service); err != nil {
    return err, ""
  } else {
    // If vermell is runningin standalone mode, all the sample have to arrived from check system to service system.
    // to ensure that, you could compare the checkStatusCache's length to the Checks one
    // that will work because in standalone mode the GetServicesStatusHuman is launch once all checks has been executed.
    if env.Context.ExecutionMode == "standalone" {
      obj = obj.WaitAllSamples(10)
      //env.Output.WriteChDebug(obj)
    }
    return nil, "Service '"+obj.GetName()+"' status is " + sample.Itoa(obj.GetStatus())
  }
}
//
//# GetServiceExitStatus: return the status for a service object to string
func (sys *ServiceSystem) GetServiceStatus(service string) (error , int) {
  ss := sys.GetServices()
  if err, obj := ss.GetServiceObject(service); err != nil {
    return err, -1
  } else {
    return nil, obj.GetStatus()    
  }
}
//
//# AddService: method add a new service to be checked
func (s *ServiceSystem) AddServiceObject(obj *ServiceObject) error {
  if err := s.Ss.AddServiceObject(obj); err != nil {
    return err
  }
  return nil
}
//
//# GetService: method returns a ServiceObject
func (s *ServiceSystem) GetServiceObject(name string) (error, *ServiceObject){
  return s.Ss.GetServiceObject(name)
}
//
//# GetServiceForCheck: method returns the services that a check is defined to
func (s *ServiceSystem) GetServicesForCheck(check string) (error, []string) {
  return s.Ss.GetServicesForCheck(check)
}

//#
//# Common methods
//#---------------------------------------------------------------------

//
//# String: converts a SampleSystem object to string
func (sys *ServiceSystem) String() string {
  return utils.ObjectToJsonString(sys.GetServices())
}

//#######################################################################################################
