package api

import (
	"verdmell/ui"
)

//#
//# Route methods
//#---------------------------------------------------------------------
//#
//# Methods for service router

//
//# GenerateAPIRoutesForService: generate a set of routes to serve
func (a* ApiEngine) GenerateAPIRoutesForService() {
	env.Output.WriteChDebug("(ApiEngine::GenerateAPIRoutesForService)")
	a.AddRoute(ui.GenerateRoute("allservices","GET","/api/services",GetAllServices))
	a.AddRoute(ui.GenerateRoute("service","GET","/api/services/{service}",GetService))
}

//#######################################################################################################