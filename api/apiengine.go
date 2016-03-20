package api

import (
	"verdmell/ui"
	"verdmell/environment"
	"verdmell/check"
	"verdmell/sample"
	"verdmell/service"
	"net/http"
)

//
var env *environment.Environment

type ApiEngine struct {
	Routes []*ui.Route
}

//#
//#
//# ApiEngine struct:
//# ApiEngine defines all informataion required for API
func NewApiEngine(e *environment.Environment) *ApiEngine {	
  e.Output.WriteChDebug("(ApiEngine::NewApiEngine)")
  sys := new(ApiEngine)
  
  // set the environment attribute
  env = e

  sys.GenerateAPIRoutes()
  e.Output.WriteChInfo("(ApiEngine::NewApiEngine) Hi! I'm your new API engine instance")
  return sys
}
//
//# SetRoutes: set routes for API
func (a *ApiEngine) SetRoutes(routes []*ui.Route){
	env.Output.WriteChDebug("(ApiEngine::SetRoutes)")
	a.Routes = routes
}
//
//# GetRoutes: set routes for API
func (a *ApiEngine) GetRoutes() []*ui.Route{
	env.Output.WriteChDebug("(ApiEngine::GetRoutes)")
	return a.Routes
}

//#
//# Specific methods
//#---------------------------------------------------------------------
//
//# SayHi: 
func (sys *ApiEngine) SayHi() {
  env.Output.WriteChInfo("(ApiEngine::SayHi) Hi! I'm your API engine instance")
}
//
//# apiWriter: write data to response writer
func apiWriter(fn func (*http.Request)(error,[]byte)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		env.Output.WriteChDebug("(ApiEngine::apiWriter) clousure")
		if err, data := fn(r); err != nil {
			http.NotFound(w,r)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}
	}
}

//
//# AddRoute: for include a new route to API Routes
func (a* ApiEngine) AddRoute(route *ui.Route){
	env.Output.WriteChDebug("(ApiEngine::AddRoute) Add new route '"+route.Name+"'")
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