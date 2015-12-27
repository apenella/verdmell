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
	a.AddRoute(ui.GenerateRoute("NodeStatus","GET","/api/node",apiWriter(GetNodeInfo)))
	a.AddRoute(ui.GenerateRoute("NodeStatus","GET","/api/node/status",apiWriter(GetNodeStatus)))
	a.AddRoute(ui.GenerateRoute("StartCheckEngine","GET","/api/node/run",apiWriter(StartCheckEngine)))
	//node checks
	a.AddRoute(ui.GenerateRoute("allchecks","GET","/api/node/checks",apiWriter(GetAllChecks)))
	a.AddRoute(ui.GenerateRoute("check","GET","/api/node/checks/{check}",apiWriter(GetCheck)))
	//node samples
	a.AddRoute(ui.GenerateRoute("allservices","GET","/api/node/samples",apiWriter(GetAllSamples)))
	a.AddRoute(ui.GenerateRoute("service","GET","/api/node/samples/{sample}",apiWriter(GetSample)))
	//node services
	a.AddRoute(ui.GenerateRoute("allservices","GET","/api/node/services",apiWriter(GetAllServices)))
	a.AddRoute(ui.GenerateRoute("service","GET","/api/node/services/{service}",apiWriter(GetService)))
}