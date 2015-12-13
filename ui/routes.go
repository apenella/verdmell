package ui

import(
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}


//#
//# Specific methods
//#---------------------------------------------------------------------

func GenerateRoute(n string, m string, p string, h http.HandlerFunc) *Route {
	route := &Route{
		Name: n,
		Method: m,
		Pattern: p,
		HandlerFunc: h,
	}
	return route
}

//#######################################################################################################