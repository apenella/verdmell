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
func (a* ApiEngine) GenerateAPIRoutesForCheck() {
	env.Output.WriteChDebug("(ApiEngine::GenerateAPIRoutesForCheck)")
	a.AddRoute(ui.GenerateRoute("allchecks","GET","/api/checks",apiWriter(GetAllChecks)))
	a.AddRoute(ui.GenerateRoute("check","GET","/api/checks/{check}",apiWriter(GetCheck)))
}

//#######################################################################################################