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
//# GenerateAPIRoutesForNode: generate a set of routes to serve
func (a* ApiEngine) GenerateAPIRoutesForNode() {
	env.Output.WriteChDebug("(ApiEngine::GenerateAPIRoutesForService)")
	a.AddRoute(ui.GenerateRoute("NodeStatus","GET","/api/node",GetNodeStatus))
	a.AddRoute(ui.GenerateRoute("StartCheckEngine","GET","/api/node/run",StartCheckEngine))
	a.AddRoute(ui.GenerateRoute("NodeStatus","GET","/api/node/status",GetNodeStatus))
}