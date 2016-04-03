/*

	Package 'ui' 
	-server
	-handler
	-router
	-routes

	-html/
	-images/
	-pages/
	-scripts/
	-style/

*/
package ui


//#
//# Route methods
//#---------------------------------------------------------------------
//#
//# Main router methods

//
//# GenerateAPIRoutes: generate a set of routes to serve
func (u* UI) GenerateRoutes() {
	env.Output.WriteChDebug("(UI::router::GenerateRoutes)")
	u.router.HandleFunc("/", u.uiHandlerFunc(Index))

	u.AddRoute(GenerateRoute("Index","GET","/", u.uiHandlerFunc(Index)))
	u.AddRoute(GenerateRoute("WebUI","GET","/ui", u.uiHandlerFunc(WebUI)))
	u.AddRoute(GenerateRoute("WebUI","GET","/sse", u.uiHandlerFunc(SSE)))
}

//
//# AddRoutes: add a set of routes 
func (u *UI) AddRoutes(routes []*Route){
	env.Output.WriteChDebug("(UI::router::AddRoutes)")
	for _,route := range routes {
		u.AddRoute(route)
	}
}

//
//# AddRoute: add one route to the router
func (u *UI) AddRoute(route *Route) {
	env.Output.WriteChDebug("(UI::router::AddRoute) Route name '"+route.Name+"'")
	u.router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(route.HandlerFunc)
}

//#######################################################################################################