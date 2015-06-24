package ui

import(
	"log"
	"net/http"

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
	u.GenerateRoutes()
	log.Fatal(http.ListenAndServe(u.listenaddr, u.router))
}

//#######################################################################################################