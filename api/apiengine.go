package api

import (
	"verdmell/ui"
	"verdmell/environment"
	"verdmell/check"
	"verdmell/sample"
	"verdmell/service"
)

//
var env *environment.Environment
var box *ObjectsBox

// constants
const(
	CHECKS = "checks"
	SERVICES = "services"
	SAMPLES = "samples"
)

type ApiEngine struct {
	Box *ObjectsBox
	Routes []*ui.Route
}

//#
//#
//# ApiEngine struct:
//# ApiEngine defines all informataion required for API
func NewApiEngine(e *environment.Environment, b *ObjectsBox) *ApiEngine {	
  e.Output.WriteChDebug("(ApiEngine::NewApiEngine)")
  sys := new(ApiEngine)
  
  // set the environment attribute
  env = e
  box = b

  sys.GenerateAPIRoutes()

  return sys
}

//# SetObjectBox: set objects for API
// func (a *ApiEngine) SetObjectBox(box *ObjectsBox){
// 	env.Output.WriteChDebug("(ApiEngine::SetObjectBox)")
// 	a.Box = box
// }
//# SetRoutes: set routes for API
func (a *ApiEngine) SetRoutes(routes []*ui.Route){
	env.Output.WriteChDebug("(ApiEngine::SetRoutes)")
	a.Routes = routes
}

//# GetObjectBox: get objects from API
// func (a *ApiEngine) GetObjectBox() *ObjectsBox {
// 	env.Output.WriteChDebug("(ApiEngine::GetObjectBox)")
// 	return a.Box
// }
//# GetRoutes: set routes for API
func (a *ApiEngine) GetRoutes() []*ui.Route{
	env.Output.WriteChDebug("(ApiEngine::GetRoutes)")
	return a.Routes
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# AddRoute: for include a new route to API Routes
func (a* ApiEngine) AddRoute(route *ui.Route){
	env.Output.WriteChDebug("(ApiEngine::AddRoute) Add new route")
	a.Routes = append(a.Routes, route)
}
//
//# GetCheckEngine: return CHECKS from obect box
func (a *ApiEngine) GetCheckEngine() *check.CheckEngine {
	if obj := box.GetObject(CHECKS); obj != nil{
		return obj.(*check.CheckEngine)
	}
	env.Output.WriteChDebug("(ApiEngine::GetCheckEngine) There is no object for "+CHECKS)
	return nil
}
//
//# GetServiceEngine: return CHECKS from obect box
func (a *ApiEngine) GetServiceEngine() *service.ServiceEngine {
	if obj := box.GetObject(SERVICES); obj != nil{
		return obj.(*service.ServiceEngine)
	}
	env.Output.WriteChDebug("(ApiEngine::GetServiceEngine) There is no object for "+SERVICES)
	return nil
}
//
//# GetSampleEngine: return CHECKS from obect box
func (a *ApiEngine) GetSampleEngine() *sample.SampleEngine {
	if obj := box.GetObject(SAMPLES); obj != nil{
		return obj.(*sample.SampleEngine)
	}
	env.Output.WriteChDebug("(ApiEngine::GetSampleEngine) There is no object for "+SAMPLES)
	return nil
}

//#######################################################################################################