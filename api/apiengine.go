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
	env.Output.WriteChDebug("(ApiEngine::GetCheckEngine)")
	return env.GetCheckEngine().(*check.CheckEngine)
}
//
//# GetServiceEngine: return CHECKS from obect box
func (a *ApiEngine) GetServiceEngine() *service.ServiceEngine {
	env.Output.WriteChDebug("(ApiEngine::GetServiceEngine)")
	return env.GetServiceEngine().(*service.ServiceEngine)
}
//
//# GetSampleEngine: return CHECKS from obect box
func (a *ApiEngine) GetSampleEngine() *sample.SampleEngine {
	env.Output.WriteChDebug("(ApiEngine::GetSampleEngine)")
	return env.GetSampleEngine().(*sample.SampleEngine)
}

//#######################################################################################################