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
	"verdmell/environment"
	"verdmell/utils"
	"verdmell/sample"
)

//
var env *environment.Environment

//#
//#
//# ServiceEngine struct:
//# ServiceEngine defines a map to store the maps
type ServiceEngine struct{
	Ss *Services `json:"servicesroot"`
	inputSampleChan chan *sample.CheckSample `json:"-"`
}

//
//# NewCheckSystem: return a Checksystem instance to be run
func NewServiceEngine(e *environment.Environment) (error, *ServiceEngine){
	e.Output.WriteChDebug("(ServiceEngine::NewServiceEngine)")
	sys := new(ServiceEngine)
	
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
	sys.StartServiceEngineReceiver()

	return nil, sys
}

//#
//# Getters/Setters methods for ServiceEngine object
//#---------------------------------------------------------------------

//
//# SetServices: methods sets the Services' value
func (s *ServiceEngine) SetServices(ss *Services) {
	s.Ss = ss
}
//
//# SetInputSampleChan: methods sets the inputSampleChan's value
func (s *ServiceEngine) SetInputSampleChan(c chan *sample.CheckSample) {
	s.inputSampleChan = c
}

//
//# GetServices: methods gets the Services' value
func (s *ServiceEngine) GetServices() *Services {
	return s.Ss
}
//
//# GetInputSampleChan: methods sets the inputSampleChan's value
func (s *ServiceEngine) GetInputSampleChan() chan *sample.CheckSample {
	return s.inputSampleChan
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# StartServiceEngine: method prepares the system to wait sample and calculate the results for services
func (s *ServiceEngine) StartServiceEngineReceiver() error {
	s.inputSampleChan = make(chan *sample.CheckSample)
	services := s.GetServices()

	env.Output.WriteChDebug("(ServiceEngine::StartServiceEngineReceiver) Starting sample receiver")
	go func() {
		defer close (s.inputSampleChan)
		for{
			select{
			case sample := <-s.inputSampleChan:
				env.Output.WriteChDebug("(ServiceEngine::StartServiceEngineReceiver) New sample received for '"+sample.GetCheck()+"'")
				_,servicesCheck := s.GetServicesForCheck(sample.GetCheck())
				for _,service := range servicesCheck {
					_,srv := services.GetServiceObject(service)
					env.Output.WriteChDebug("(ServiceEngine::StartServiceEngineReceiver) Sample for '"+sample.GetCheck()+"' belongs to '"+srv.GetName()+"'")
					go srv.SendToSampleChannel(sample)
				}
			}
		}
	}()
	return nil
}
//
//# SendSampleToServiceEngine: method prepares the system to wait sample and calculate the results for services
func (s *ServiceEngine) SendSampleToServiceEngine(sample *sample.CheckSample) {
	env.Output.WriteChDebug("(ServiceEngine::SendSampleToServiceEngine) Send sample "+sample.String())
	s.inputSampleChan <- sample
}
//
//# RegisterService: register a new service for ServiceSysem
func (s *ServiceEngine) RegisterService(name string, desc string, checks []string) error {
	env.Output.WriteChDebug("(ServiceEngine::RegisterService) New service to register '"+name+"'")
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
//# GetAllServices: return information for all services
func (sys *ServiceEngine) GetAllServices() []byte {
	env.Output.WriteChDebug("(ServiceEngine::GetAllServices)")
	ss := sys.GetServices()

	//return ss.String()
	return utils.ObjectToJsonByte(ss)
}
//
//# GetServices: return all information about a service
func (sys *ServiceEngine) GetService(service string) []byte {
	env.Output.WriteChDebug("(ServiceEngine::GetService)")
	ss := sys.GetServices()
	s := ss.GetServices()

	//return s[service].String()
	return utils.ObjectToJsonByte(s[service])
}



//
//# GetAllServicesStatusHuman: converts a SampleSystem object to string
func (sys *ServiceEngine) GetAllServicesStatusHuman() (error, string) {
	env.Output.WriteChDebug("(ServiceEngine::GetAllServicesStatusHuman)")
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
func (sys *ServiceEngine) GetServicesStatusHuman(service string) (error ,string) {
	env.Output.WriteChDebug("(ServiceEngine::GetServicesStatusHuman) Get status for '"+service+"'")
	var obj *ServiceObject
	var err error
	srvChan := make(chan *ServiceObject)
	defer close(srvChan)

	ss := sys.GetServices()
	if err, obj = ss.GetServiceObject(service); err != nil {
		return err, ""
	} else {
		// If vermell is running as standalone mode, all the sample have to arrived from check system to service system.
		// to ensure that, you could compare the checkStatusCache's length to the Checks one
		// that will work because in standalone mode the GetServicesStatusHuman is launch once all checks has been executed.
		if env.Context.ExecutionMode == "standalone" {
			obj = obj.WaitAllSamples(5)
			env.Output.WriteChDebug("(ServiceEngine::GetServicesStatusHuman) The waiting has end")
		}
		return nil, "Service '"+obj.GetName()+"' status is " + sample.Itoa(obj.GetStatus())
	}
}
//
//# GetServiceExitStatus: return the status for a service object to string
func (sys *ServiceEngine) GetServiceStatus(service string) (error , int) {
	ss := sys.GetServices()
	if err, obj := ss.GetServiceObject(service); err != nil {
		return err, -1
	} else {
		// If vermell is running as standalone mode, all the sample have to arrived from check system to service system.
		// to ensure that, you could compare the checkStatusCache's length to the Checks one
		// that will work because in standalone mode the GetServicesStatusHuman is launch once all checks has been executed.
		if env.Context.ExecutionMode == "standalone" {
			obj = obj.WaitAllSamples(5)
			env.Output.WriteChDebug("(ServiceEngine::GetServicesStatusHuman) The waiting has end")
		}
		return nil, obj.GetStatus()		
	}
}
//
//# AddService: method add a new service to be checked
func (s *ServiceEngine) AddServiceObject(obj *ServiceObject) error {
	if err := s.Ss.AddServiceObject(obj); err != nil {
		return err
	}
	return nil
}
//
//# GetService: method returns a ServiceObject
func (s *ServiceEngine) GetServiceObject(name string) (error, *ServiceObject){
	return s.Ss.GetServiceObject(name)
}
//
//# GetServiceForCheck: method returns the services that a check is defined to
func (s *ServiceEngine) GetServicesForCheck(check string) (error, []string) {
	return s.Ss.GetServicesForCheck(check)
}

//#
//# Common methods
//#---------------------------------------------------------------------

//
//# String: converts a SampleSystem object to string
func (sys *ServiceEngine) String() string {
	return utils.ObjectToJsonString(sys.GetServices())
}

//#######################################################################################################
