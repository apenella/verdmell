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
func (a* ApiSystem) GenerateAPIRoutesForService() {
	env.Output.WriteChDebug("(ApiSystem::GenerateAPIRoutesForService)")
	a.AddRoute(ui.GenerateRoute("allservices","GET","/api/services",GetAllServices))
	a.AddRoute(ui.GenerateRoute("service","GET","/api/services/{service}",GetService))
}

//#######################################################################################################