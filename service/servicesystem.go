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
  Ss *Services `json: ""`
  inputSampleChan chan *sample.CheckSample `json:"-"`
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
//# GetServiceForCheck: method returns the services that a check is defined to
func (s *ServiceSystem) GetServicesForCheck(check string) (error, []string) {
  var err error
  var services []string

  ss := s.GetServices()
  if err, services = ss.GetServicesForCheck(check); err != nil {
    return err, nil
  }

  return nil, services
}

//
//# ServicesStatusHuman: converts a SampleSystem object to string
func (sys *ServiceSystem) ServicesStatusHuman() string {
  var str string
  ss := sys.GetServices()

  for _,obj := range ss.GetServices(){
    str += "Service '"+obj.GetName()+"' status is " + sample.Itoa(obj.GetStatus()) + "\n"
  }

  return str
}

//#
//# Common methods
//#---------------------------------------------------------------------

//
//# String: converts a SampleSystem object to string
func (sys *ServiceSystem) String() string {
  return utils.ObjectToJsonString(sys)
}


//#######################################################################################################
