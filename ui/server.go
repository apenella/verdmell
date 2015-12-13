package ui

import(
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

var ui *UI = nil

type UI struct{
	listenaddr string
	router *mux.Router
}

func NewUI(listenaddr string) *UI {
	if ui == nil {
		ui = &UI{
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
	u.router.Handle("/images/{img}",http.StripPrefix("/images/", http.FileServer(http.Dir("./ui/images/"))))
	log.Fatal(http.ListenAndServe(u.listenaddr, u.router))
}

//#######################################################################################################