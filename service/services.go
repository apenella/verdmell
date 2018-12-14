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
	"os"
	"verdmell/utils"
)

//#
//#
//# ServiceEngine struct:
//# ServiceEngine defines a map to store the maps
type Services struct {
	Services              map[string]*ServiceObject `json:"services"`
	checkServiceMapReduce map[string][]string       `json:"-"`
}

//#
//# Getters/Setters methods for Checks object
//#---------------------------------------------------------------------

//# SetService: methods sets the Services value for the Services object
func (ss *Services) SetServices(s map[string]*ServiceObject) {
	ss.Services = s
}

//
//# SetCheckServiceMapReduce: methods sets the CheckServiceMapReduce's value
func (s *Services) SetCheckServiceMapReduce(mr map[string][]string) {
	s.checkServiceMapReduce = mr
}

//# GetService: methods gets the Services's value for a gived Services object
func (ss *Services) GetServices() map[string]*ServiceObject {
	return ss.Services
}

//
//# GetCheckServiceMapReduce: methods sets the CheckServiceMapReduce's value
func (s *Services) GetCheckServiceMapReduce() map[string][]string {
	return s.checkServiceMapReduce
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# AddServiceObject: method add a new service to be checked
func (s *Services) AddServiceObject(obj *ServiceObject) error {
	name := obj.GetName()
	if _, exist := s.Services[name]; !exist {
		s.Services[name] = obj
	} else {
		env.Output.WriteChWarn("(Services::AddService) The service '" + name + "' is already defined")
	}
	return nil
}

//
//# GetServiceObject: method returns a ServiceObject
func (s *Services) GetServiceObject(name string) (error, *ServiceObject) {
	var exist bool
	var srv *ServiceObject

	if srv, exist = s.Services[name]; !exist {
		err := errors.New("(Services::GetService) Service '" + name + "' doesn't exist.")
		return err, nil
	}

	return nil, srv
}

//
//# GetServiceForCheck: method returns the services that a check is defined to
func (s *Services) GetServicesForCheck(check string) (error, []string) {
	var exist bool
	var services []string

	if services, exist = s.checkServiceMapReduce[check]; !exist {
		err := errors.New("(Services::GetServicesForCheck) Check '" + check + "' has any service defined.")
		return err, nil
	}

	return nil, services
}

//
//# UnmarshalService: get the json content from a file and field an Services object on it.
//	The method requieres a file path.
//	The method returns a pointer to Services object
func UnmarshalServices(file string) *Services {
	env.Output.WriteChDebug("(Services::UnmarshalServices)")

	s := new(Services)
	// extract the content from the file and dumps it on the CHecks object
	utils.LoadJSONFile(file, s)

	return s
}

//
//# RetrieveChecks: gets all the files found on checks folder and generate one Checks object with all this CheckObject defined.
func RetrieveServices(folder string) *Services {
	services := new(Services)
	// checks will contain all the CheckObject definition
	servicesMap := make(map[string]*ServiceObject)
	// files is an array with all files found inside the folder
	files := utils.GetFolderFiles(folder)

	// sync channel
	serviceObjChan := make(chan *ServiceObject)
	serviceskFileEndChan := make(chan bool)
	allServicesGetChan := make(chan bool)
	done := make(chan *Services)

	// goroutine for extract each check object from file
	retrieveServicesFromFile := func(f os.FileInfo) {
		serviceFile := folder + string(os.PathSeparator) + f.Name()
		env.Output.WriteChDebug("(Services::RetrieveServices) File found: " + serviceFile)

		s := UnmarshalServices(serviceFile)

		if len(s.GetServices()) == 0 {
			env.Output.WriteChWarn("(Services::RetrieveServices) You should review the file " + serviceFile + ", no service has been load from it")
		}

		for servicekName, servicekObj := range s.GetServices() {

			servicekObj.SetName(servicekName)
			servicekObj.SetStatus(-1)
			// sending the CheckObject to be stored
			serviceObjChan <- servicekObj
			env.Output.WriteChInfo("(Services::RetrieveServices) Service '" + servicekName + "' defined")
			env.Output.WriteChDebug("(Services::RetrieveServices) '" + servicekObj.String() + "'")
		}
		// a message is send when all ServiceObject defined into a file have been sent to store
		serviceskFileEndChan <- true
	}
	// call the goroutine for each file
	for _, f := range files {
		go retrieveServicesFromFile(f)
	}
	// waiting for all serviceskFileEndChan that will indicate that all files has been analized
	go func() {
		for i := len(files); i > 0; i-- {
			<-serviceskFileEndChan
		}
		defer close(serviceskFileEndChan)
		allServicesGetChan <- true
	}()
	// store all ServiceObjects sent. Once the allChecksGetChan channel gets a message the goroutine will assume that all ServicesOjects has been sent
	go func() {
		var services Services
		allServicesGet := false
		for !allServicesGet {
			select {
			// get a CheckObject object
			case srv := <-serviceObjChan:
				env.Output.WriteChDebug("(Services::RetrieveServices::routine) New service to register '" + srv.GetName() + "'")
				if _, exist := servicesMap[srv.GetName()]; !exist {
					servicesMap[srv.GetName()] = srv
					go func() {
						srv.StartReceiver()
					}()
				}
			// ending message
			case allServicesGet = <-allServicesGetChan:
				services.SetServices(servicesMap)
				done <- &services
				defer close(serviceObjChan)
				defer close(allServicesGetChan)
			}
		}
	}()
	// the main routine will wait for the work to be done
	services = <-done
	defer close(done)

	return services
}

//
//# ValidateServices: ensures that all services are defined properly
func (s *Services) ValidateServices() error {
	errorChan := make(chan error)
	statusChan := make(chan bool)

	//goroutine to validate each service
	validation := func(obj *ServiceObject) {
		env.Output.WriteChDebug("(Services::ValidateServices::routine) Validating services '" + obj.GetName() + "'")
		if err := obj.ValidateServiceObject(); err != nil {
			errorChan <- err
		} else {
			statusChan <- true
		}
	}

	// for each ServiceObject is launched a validation function
	for _, srv := range s.GetServices() {
		go validation(srv)
	}

	// the method waits for all the status. If an error occurs, the function returns it
	for i := 0; i < len(s.GetServices()); i++ {
		select {
		case err := <-errorChan:
			close(errorChan)
			return err
		case <-statusChan:
			break
		}
	}

	close(statusChan)
	// if no error has been found, all ServiceObjects have been defined correctly
	return nil
}

//
//# GenerateCheckServices: method return the relationship between all checks and its services
func (s *Services) GenerateCheckServices() error {
	//based on mapReduce
	mapChan := make(chan map[string]string)
	defer close(mapChan)
	reduceChan := make(chan map[string][]string)
	defer close(reduceChan)

	// The Map
	// map generate a map for each service containing a relationship between the checks and the service
	generateChecksServiceMap := func(service string, serviceObj *ServiceObject) {
		// map to store
		csMap := make(map[string]string)
		env.Output.WriteChDebug("(Services::GenerateCheckServices::generateChecksServiceMap) for service '" + service + "'")
		for _, check := range serviceObj.GetChecks() {
			env.Output.WriteChDebug("(Services::GenerateCheckServices::generateChecksServiceMap) The check '" + check + "' used for service '" + service + "'")
			csMap[check] = service
		}
		mapChan <- csMap
	}
	// The Reduce
	// merge all the maps results
	generateChecksServiceReduce := func() {
		csReduce := make(map[string][]string)

		for i := 0; i < len(s.GetServices()); i++ {
			select {
			case csMap := <-mapChan:
				for check, service := range csMap {
					env.Output.WriteChDebug("(Services::GenerateCheckServices::generateChecksServiceReduce) The '" + check + "' used for service '" + service + "'")
					csReduce[check] = append(csReduce[check], service)
				}
			}
		}
		reduceChan <- csReduce
	}

	for service, serviceObj := range s.GetServices() {
		go generateChecksServiceMap(service, serviceObj)
	}

	go generateChecksServiceReduce()

	cs := <-reduceChan
	s.SetCheckServiceMapReduce(cs)

	return nil
}

//
// String method converts a Services object to string
func (s *Services) String() string {
	var str string
	var err error

	str, err = utils.ObjectToJSONString(s)
	if err != nil {
		return err.Error()
	}

	return str
}
