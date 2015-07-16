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
func (a* ApiSystem) GenerateAPIRoutes() {
	env.Output.WriteChDebug("(ApiSystem::GenerateAPIRoutes)")
	a.AddRoute(ui.GenerateRoute("Index","GET","/api",Index))
	a.AddRoute(ui.GenerateRoute("Cluster","GET","/api/cluster",GetCluster))

	a.GenerateAPIRoutesForCheck()
	a.GenerateAPIRoutesForService()
	a.GenerateAPIRoutesForSamples()
	a.GenerateAPIRoutesForNode()
}

//#######################################################################################################