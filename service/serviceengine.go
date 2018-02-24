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
	"errors"
	"strconv"
	
	"verdmell/environment"
	"verdmell/check"
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
	inputChannel chan interface{} `json:"-"`
	outputChannels map[chan interface{}] string `json: "-"`
}

//
//# NewCheckSystem: return a Checksystem instance to be run
func NewServiceEngine(e *environment.Environment) (error, *ServiceEngine){
	e.Output.WriteChDebug("(ServiceEngine::NewServiceEngine)")
	sys := new(ServiceEngine)
	
	// set the environment attribute
	env = e

	// folder that contains services definitions
	folder := env.Config.Services.Folder
	// Get defined services
	env.Output.WriteChDebug("(ServiceEngine::NewServiceEngine) RetrieveServices")
	srv := RetrieveServices(folder)

	// set the attribute CheckServiceMapReduce
	env.Output.WriteChDebug("(ServiceEngine::NewServiceEngine) GenerateCheckServices")
	srv.GenerateCheckServices()
	
	env.Output.WriteChDebug("(ServiceEngine::NewServiceEngine) ValidateServices")
	if err := srv.ValidateServices(); err != nil {
		return err, nil
	}
	env.Output.WriteChDebug("(ServiceEngine::NewServiceEngine) SetServices")
	sys.SetServices(srv)

	// Set description for default service

	desc := "Global services for node "+env.Config.Name
	
	env.Output.WriteChDebug("(ServiceEngine::NewServiceEngine) Registering service '"+env.Config.Name+"'")
	checkengine := env.GetCheckEngine().(*check.CheckEngine)
	if err := sys.RegisterService(env.Config.Name,desc, checkengine.ListCheckNames()); err != nil{
		return err, nil
	}

	// Initialize the OutputChannels
  sys.outputChannels = make(map[chan interface{}] string)

	// start the sample receiver
	env.Output.WriteChDebug("(ServiceEngine::NewServiceEngine) Start")
	sys.Start()

	// Set the environments services engine
	env.SetServiceEngine(sys)
	env.Output.WriteChInfo("(ServiceEngine::NewServiceEngine) Hi! I'm your new service engine instance")

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
//# SetinputChannel: methods sets the inputChannel's value
func (s *ServiceEngine) SetInputChannel(c chan interface{}) {
	s.inputChannel = c
}
//
//# SetOutputChannels: method sets the channels to write service status
func (s *ServiceEngine) SetOutputChannels(o map[chan interface{}] string) {
  env.Output.WriteChDebug("(ServiceEngine::SetOutputChannels)")
  s.outputChannels = o
}

//
//# GetServices: methods gets the Services' value
func (s *ServiceEngine) GetServices() *Services {
	return s.Ss
}
//
//# GetinputChannel: methods sets the inputChannel's value
func (s *ServiceEngine) GetInputChannel() chan interface{} {
	return s.inputChannel
}
//
//# GetOutputChannels: methods return the channels to write samples
func (s *ServiceEngine) GetOutputChannels() map[chan interface{}] string {
  env.Output.WriteChDebug("(ServiceEngine::GetOutputChannels)")
  return s.outputChannels
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# SayHi: 
func (s *ServiceEngine) SayHi() {
  env.Output.WriteChInfo("(ServiceEngine::SayHi) Hi! I'm your new service engine instance")
}
//
//# Subscribe: Add a new channel to write service status
func (s *ServiceEngine) Subscribe(o chan interface{}, desc string) error {
  env.Output.WriteChDebug("(ServiceEngine::Subscribe)")

  channels := s.GetOutputChannels()
  if _, exist := channels[o]; !exist {
    channels[o] = desc
  } else {
    return errors.New("(ServiceEngine::Subscribe) You are trying to add an existing channel")
  }

  return nil
}

//
//# StartServiceEngine: method prepares the system to wait sample and calculate the results for services
func (s *ServiceEngine) Start() error {
	s.inputChannel = make(chan interface{})
	services := s.GetServices()

	env.Output.WriteChDebug("(ServiceEngine::Start) Starting sample receiver")
	go func() {
		defer close (s.inputChannel)
		for{
			select{
			case obj := <-s.inputChannel:
				sample := obj.(*sample.CheckSample)
				env.Output.WriteChDebug("(ServiceEngine::Start) New sample received for '"+sample.GetCheck()+"'")
				_,servicesCheck := s.GetServicesForCheck(sample.GetCheck())
				for _,service := range servicesCheck {
					_,srv := services.GetServiceObject(service)
					env.Output.WriteChDebug("(ServiceEngine::StartReceiver) Sample for '"+sample.GetCheck()+"' belongs to '"+srv.GetName()+"'")
					go srv.RecevieData(sample)
				}
			}
		}
	}()
	return nil
}

//
//# SendData: method that send services to other engines
func (s *ServiceEngine) SendData(o *ServiceObject) error {
	env.Output.WriteChDebug("(ServiceEngine::SendData)")
	
	for c,desc := range s.GetOutputChannels(){
		env.Output.WriteChDebug("(ServiceEngine::SendData) Writing service to channel '"+desc+"' {service:'"+o.GetName()+"', status:"+strconv.Itoa(o.GetStatus())+", timestamp:"+strconv.Itoa(int(o.GetTimestamp()))+"}")
		c <- o
	}

	return nil
}
//
//# ReceiveData: method prepares the system to wait sample and calculate the results for services
//func (s *ServiceEngine) SendSample(sample *sample.CheckSample) {
func (s *ServiceEngine) ReceiveData(sample *sample.CheckSample) {
	env.Output.WriteChDebug("(ServiceEngine::ReceiveData) Send sample "+sample.String())
	s.inputChannel <- sample
}

//
//# RegisterService: register a new service for service engine
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

	env.Output.WriteChDebug("(ServiceEngine::RegisterService) Service '"+name+"' registered properly")

	return nil
}


//
//# GetAllServices: return information for all services
func (sys *ServiceEngine) GetAllServices() (error, []byte) {
	env.Output.WriteChDebug("(ServiceEngine::GetAllServices)")
	var services *Services

	if services = sys.GetServices(); services== nil {
		msg := "(ServiceEngine::GetService) There are no services defined."
		env.Output.WriteChDebug(msg)
		return errors.New(msg), nil
	}

	//return ss.String()
	return utils.ObjectToJsonByte(services)
}
//
//# GetServices: return all information about a service
func (sys *ServiceEngine) GetService(name string) (error, []byte) {
	env.Output.WriteChDebug("(ServiceEngine::GetService)")
	var services *Services
	var service map[string] *ServiceObject
	var obj *ServiceObject
	var exist bool

	if services = sys.GetServices(); services == nil {
		msg := "(ServiceEngine::GetService) There are no services defined."
		env.Output.WriteChDebug(msg)
		return errors.New(msg), nil
	}

	if service = services.GetServices(); service == nil {
		msg := "(ServiceEngine::GetService) There are no services defined."
		env.Output.WriteChDebug(msg)
		return errors.New(msg), nil
	}

	if obj, exist = service[name]; !exist {
		msg := "(ServiceEngine::GetService) The service '"+name+"' is not defined."
		env.Output.WriteChDebug(msg)
		return errors.New(msg), nil
	}

	return utils.ObjectToJsonByte(obj)
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
		//if env.Context.ExecutionMode == "standalone" {
			obj = obj.WaitAllSamples(5)
			env.Output.WriteChDebug("(ServiceEngine::GetServicesStatusHuman) The waiting has end")
		//}
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
		//if env.Context.ExecutionMode == "standalone" {
			obj = obj.WaitAllSamples(5)
			env.Output.WriteChDebug("(ServiceEngine::GetServiceStatus) The waiting has end")
		//}
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
	if err, str := utils.ObjectToJsonString(sys.GetServices()); err != nil{
    return err.Error()
  } else{
    return str
  }
}

//#######################################################################################################
