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
func (a* ApiSystem) GenerateAPIRoutesForCheck() {
	env.Output.WriteChDebug("(ApiSystem::GenerateAPIRoutesForCheck)")
	a.AddRoute(ui.GenerateRoute("allchecks","GET","/api/checks",GetAllChecks))
	a.AddRoute(ui.GenerateRoute("check","GET","/api/checks/{check}",GetCheck))
}

//#######################################################################################################