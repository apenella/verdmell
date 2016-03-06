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

import(
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"path"
	"verdmell/environment"
)

//
var env *environment.Environment
var ui *UI = nil


type UI struct{
	Listenaddr string
	router *mux.Router
	templates *template.Template
}
//
//# NewUI: return a new UI
func NewUI(e *environment.Environment, listenaddr string) *UI {
	// if it's already running an UI instance is not created a new one
	if ui == nil {
		env = e
		index := path.Join("ui","html", "index.html")
		header := path.Join("ui","html", "header.html")
		content := path.Join("ui","html", "content.html")
		footer := path.Join("ui","html", "footer.html")
		//scripts := path.Join("ui","scripts", "scripts.js")
		scripts := path.Join("ui","scripts", "verdmell.js")
		style := path.Join("ui","style", "verdmell.css")

		ui = &UI{
			Listenaddr: listenaddr,
			router: mux.NewRouter().StrictSlash(true),
			templates: template.Must(template.ParseFiles(index,scripts,style,header,content,footer)),
		}
	}
	env.Output.WriteChDebug("(UI::server::NewUI) New UI listening at: "+ui.Listenaddr)
	return ui
}


func GetUI() *UI {
	env.Output.WriteChDebug("(UI::server::GetUI) Get UI listening at: "+ui.Listenaddr)
	return ui
}

func (u *UI) StartUI(){
	env.Output.WriteChDebug("(UI::server::StartUI) Starting UI listening at: "+u.Listenaddr)
	u.GenerateRoutes()
	u.router.Handle("/images/{img}",http.StripPrefix("/images/", http.FileServer(http.Dir("./ui/images/"))))
	u.router.Handle("/scripts/{script}",http.StripPrefix("/scripts/", http.FileServer(http.Dir("./ui/scripts/"))))
	u.router.Handle("/style/{style}",http.StripPrefix("/style/", http.FileServer(http.Dir("./ui/style/"))))
	log.Fatal(http.ListenAndServe(u.Listenaddr, u.router))
}

//#
//# Specific methods
//#---------------------------------------------------------------------


//
//# apiWriter: write data to response writer
func (u *UI) uiHandlerFunc(fn func (http.ResponseWriter,*http.Request,*UI)(error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		if err := fn(w,r,u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

//#
//# Common methods
//#---------------------------------------------------------------------

//
//# String: converts a SampleSystem object to string
func (u *UI) String() string {
  return "{ listenaddr: '"+u.Listenaddr+"' }"
}

//#######################################################################################################