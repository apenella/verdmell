package api

import (
	"verdmell/ui"
)

//#
//# Route methods
//#---------------------------------------------------------------------
//#
//# Methods for sample router

//
//# GenerateAPIRoutesForSamples: generate a set of routes to serve
func (a* ApiEngine) GenerateAPIRoutesForSamples() {
	env.Output.WriteChDebug("(ApiEngine::GenerateAPIRoutesForSamples)")
	a.AddRoute(ui.GenerateRoute("allservices","GET","/api/samples",GetAllSamples))
	a.AddRoute(ui.GenerateRoute("service","GET","/api/samples/{sample}",GetSample))
}

//#######################################################################################################