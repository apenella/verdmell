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
	listenaddr string
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
		scripts := path.Join("ui","scripts", "scripts.js")
		style := path.Join("ui","style", "verdmell.css")
		header := path.Join("ui","html", "header.html")
		content := path.Join("ui","html", "content.html")
		footer := path.Join("ui","html", "footer.html")

		ui = &UI{
			listenaddr: listenaddr,
			router: mux.NewRouter().StrictSlash(true),
			templates: template.Must(template.ParseFiles(index,scripts,style,header,content,footer)),
		}
	}
	env.Output.WriteChDebug("(UI::server::NewUI) New UI listening at: "+ui.listenaddr)
	return ui
}


func GetUI() *UI {
	env.Output.WriteChDebug("(UI::server::GetUI) Get UI listening at: "+ui.listenaddr)
	return ui
}

func (u *UI) StartUI(){
	env.Output.WriteChDebug("(UI::server::StartUI) Starting UI listening at: "+ui.listenaddr)
	u.GenerateRoutes()
	u.router.Handle("/images/{img}",http.StripPrefix("/images/", http.FileServer(http.Dir("./ui/images/"))))
	log.Fatal(http.ListenAndServe(u.listenaddr, u.router))
}

//#
//# Specific methods
//#---------------------------------------------------------------------


//
//# apiWriter: write data to response writer
func uiWriter(fn func (*http.Request)(error,[]byte)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){

		if err, data := fn(r); err != nil {
			http.NotFound(w,r)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}
	}
}

//#
//# Common methods
//#---------------------------------------------------------------------

//
//# String: converts a SampleSystem object to string
func (ui *UI) String() string {
  return "{ listenaddr: '"+ui.listenaddr+"' }"
}

//#######################################################################################################