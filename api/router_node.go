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
func (a* ApiSystem) GenerateAPIRoutesForNode() {
	env.Output.WriteChDebug("(ApiSystem::GenerateAPIRoutesForService)")
	a.AddRoute(ui.GenerateRoute("NodeStatus","GET","/api/node",GetNodeStatus))
	a.AddRoute(ui.GenerateRoute("StartCheckSystem","GET","/api/node/run",StartCheckSystem))
	a.AddRoute(ui.GenerateRoute("NodeStatus","GET","/api/node/status",GetNodeStatus))
}