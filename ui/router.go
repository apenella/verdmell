package ui


func (u *UI) AddRoutes(routes []*Route){
	for _,route := range routes {
		u.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

//#######################################################################################################