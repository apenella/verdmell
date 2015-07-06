package ui


//#
//# Route methods
//#---------------------------------------------------------------------
//#
//# Main router methods

//
//# GenerateAPIRoutes: generate a set of routes to serve
func (u* UI) GenerateRoutes() {
	u.router.HandleFunc("/", Index)

	u.AddRoute(GenerateRoute("Index","GET","/",Index))
	u.AddRoute(GenerateRoute("WebUI","GET","/ui",WebUI))
}


func (u *UI) AddRoutes(routes []*Route){
	for _,route := range routes {
		u.AddRoute(route)
	}
}

func (u *UI) AddRoute(route *Route) {
	u.router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(route.HandlerFunc)
}

//#######################################################################################################