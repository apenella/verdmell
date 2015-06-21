package ui

import(
	"log"
	"fmt"
	"net/http"
	"html"

	"github.com/gorilla/mux"
)

type UI struct{
	listenaddr string
	router *mux.Router
}

func NewUI(listenaddr string) *UI {
	ui := &UI{
		listenaddr: listenaddr,
		router: mux.NewRouter().StrictSlash(true),
	}
	return ui
}

func (u *UI) StartUI(){
	u.router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(u.listenaddr, u.router))
}

func (u *UI) AddRoutes(routes []*Route){
	for _,route := range routes {
		u.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}