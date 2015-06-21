package api

import (
	"verdmell/environment"
	"verdmell/check"
	"verdmell/sample"
	"verdmell/service"
)

//
var env *environment.Environment

// constants
var CHECKS = "CHECKS"
var SERVICES = "SERVICES"
var SAMPLES = "SAMPLES"

type ApiSystem struct {
	box *ObjectsBox
}

//#
//#
//# ApiSystem struct:
//# ApiSystem defines all informataion required for API
func NewApiSystem(e *environment.Environment) *ApiSystem {
	e.Output.WriteChDebug("(ApiSystem::NewApiSystem)")
  sys := new(ApiSystem)
  
  // set the environment attribute
  env = e

  return sys
}

//# SetObjectBox: set objects for API
func (a *ApiSystem) SetObjectBox(box *ObjectsBox){
	env.Output.WriteChDebug("(ApiSystem::SetObjectBox)")
	a.box = box
}

//# GetObjectBox: get objects from API
func (a *ApiSystem) GetObjectBox() *ObjectsBox {
	env.Output.WriteChDebug("(ApiSystem::GetObjectBox)")
	return a.box
}

//#
//# Specific methods
//#---------------------------------------------------------------------


//# GetCheckSystem: return CHECKS from obect box
func (a *ApiSystem) GetCheckSystem() *check.CheckSystem {
	if obj := a.box.GetObject(CHECKS); obj != nil{
		return obj.(*check.CheckSystem)
	}
	env.Output.WriteChDebug("(ApiSystem::GetCheckSystem) There is no object for "+CHECKS)
	return nil
}

//# GetServiceSystem: return CHECKS from obect box
func (a *ApiSystem) GetServiceSystem() *service.ServiceSystem {
	if obj := a.box.GetObject(SERVICES); obj != nil{
		return obj.(*service.ServiceSystem)
	}
	env.Output.WriteChDebug("(ApiSystem::GetServiceSystem) There is no object for "+SERVICES)
	return nil
}

//# GetSampleSystem: return CHECKS from obect box
func (a *ApiSystem) GetSampleSystem() *sample.SampleSystem {
	if obj := a.box.GetObject(SAMPLES); obj != nil{
		return obj.(*sample.SampleSystem)
	}
	env.Output.WriteChDebug("(ApiSystem::GetSampleSystem) There is no object for "+SAMPLES)
	return nil
}


