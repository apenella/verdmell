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

type ApiSystem struct {
	Box *ObjectsBox
	Routes []*ui.Route
}

//#
//#
//# ApiSystem struct:
//# ApiSystem defines all informataion required for API
func NewApiSystem(e *environment.Environment, b *ObjectsBox) *ApiSystem {
	e.Output.WriteChDebug("(ApiSystem::NewApiSystem)")
  sys := new(ApiSystem)
  
  // set the environment attribute
  env = e
  box = b

  sys.GenerateAPIRoutes()

  return sys
}

//# SetObjectBox: set objects for API
// func (a *ApiSystem) SetObjectBox(box *ObjectsBox){
// 	env.Output.WriteChDebug("(ApiSystem::SetObjectBox)")
// 	a.Box = box
// }
//# SetRoutes: set routes for API
func (a *ApiSystem) SetRoutes(routes []*ui.Route){
	env.Output.WriteChDebug("(ApiSystem::SetRoutes)")
	a.Routes = routes
}

//# GetObjectBox: get objects from API
// func (a *ApiSystem) GetObjectBox() *ObjectsBox {
// 	env.Output.WriteChDebug("(ApiSystem::GetObjectBox)")
// 	return a.Box
// }
//# GetRoutes: set routes for API
func (a *ApiSystem) GetRoutes() []*ui.Route{
	env.Output.WriteChDebug("(ApiSystem::GetRoutes)")
	return a.Routes
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# AddRoute: for include a new route to API Routes
func (a* ApiSystem) AddRoute(route *ui.Route){
	env.Output.WriteChDebug("(ApiSystem::AddRoute) Add new route")
	a.Routes = append(a.Routes, route)
}
//
//# GetCheckEngine: return CHECKS from obect box
func (a *ApiSystem) GetCheckEngine() *check.CheckEngine {
	if obj := box.GetObject(CHECKS); obj != nil{
		return obj.(*check.CheckEngine)
	}
	env.Output.WriteChDebug("(ApiSystem::GetCheckEngine) There is no object for "+CHECKS)
	return nil
}
//
//# GetServiceSystem: return CHECKS from obect box
func (a *ApiSystem) GetServiceSystem() *service.ServiceSystem {
	if obj := box.GetObject(SERVICES); obj != nil{
		return obj.(*service.ServiceSystem)
	}
	env.Output.WriteChDebug("(ApiSystem::GetServiceSystem) There is no object for "+SERVICES)
	return nil
}
//
//# GetSampleSystem: return CHECKS from obect box
func (a *ApiSystem) GetSampleSystem() *sample.SampleSystem {
	if obj := box.GetObject(SAMPLES); obj != nil{
		return obj.(*sample.SampleSystem)
	}
	env.Output.WriteChDebug("(ApiSystem::GetSampleSystem) There is no object for "+SAMPLES)
	return nil
}

//#######################################################################################################