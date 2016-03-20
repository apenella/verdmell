package api

import (
	"verdmell/ui"
)

//#
//# Route methods
//#---------------------------------------------------------------------
//#
//# Methods for check router

//
//# GenerateAPIRoutesForCheck: generate a set of routes to serve
func (a* ApiEngine) GenerateAPIRoutesForCluster() {
	env.Output.WriteChDebug("(ApiEngine::GenerateAPIRoutesForCluster)")
	a.AddRoute(ui.GenerateRoute("allchecks","GET","/api/cluster",apiWriter(GetClusterInfo)))
	a.AddRoute(ui.GenerateRoute("allchecks","GET","/api/cluster/nodes",apiWriter(GetClusterNodes)))
}

//#######################################################################################################