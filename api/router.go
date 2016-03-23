package api

import (
	"verdmell/ui"
)

//#
//# Route methods
//#---------------------------------------------------------------------
//#
//# Main router methods

//
//# GenerateAPIRoutes: generate a set of routes to serve
func (a* ApiEngine) GenerateAPIRoutes() {
	env.Output.WriteChDebug("(ApiEngine::GenerateAPIRoutes)")
	a.AddRoute(ui.GenerateRoute("Index","GET","/api",Index))
	
	a.GenerateAPIRoutesForCheck()
	a.GenerateAPIRoutesForService()
	a.GenerateAPIRoutesForSamples()
	a.GenerateAPIRoutesForNode()
	a.GenerateAPIRoutesForCluster()
}

//#######################################################################################################