package ui

import(
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

var ui *UI = nil

type UI struct{
	nodes map[string] string
	listenaddr string
	router *mux.Router
}

func NewUI(listenaddr string, n map[string] string) *UI {
	if ui == nil {
		ui = &UI{
			nodes: n,
			listenaddr: listenaddr,
			router: mux.NewRouter().StrictSlash(true),
		}
	}
	return ui
}

func GetUI() *UI {
	return ui
}

func (u *UI) StartUI(){
	u.GenerateRoutes()
	log.Fatal(http.ListenAndServe(u.listenaddr, u.router))
}



//#######################################################################################################