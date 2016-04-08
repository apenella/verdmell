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
	"errors"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
	"verdmell/environment"
)

//
var env *environment.Environment
var ui *UI = nil


type UI struct {
	Listenaddr string
	router *mux.Router
	templates *template.Template
	inputChannel chan []byte

	clients map[chan []byte]bool
	newClients chan chan []byte
	defunctClients chan chan []byte
}
//
//# NewUI: return a new UI
func NewUI(e *environment.Environment, listenaddr string) *UI {
	// if it's already running an UI instance is not created a new one

	// set environment
	env = e

	index := path.Join("ui","html", "index.html")
	header := path.Join("ui","html", "header.html")
	content := path.Join("ui","html", "content.html")
	footer := path.Join("ui","html", "footer.html")
	scripts := path.Join("ui","scripts", "verdmell.js")
	style := path.Join("ui","style", "verdmell.css")

	if ui == nil {
		ui = new(UI)
		ui.SetListenaddr(listenaddr)
		ui.SetRouter(mux.NewRouter().StrictSlash(true))
		ui.SetTemplates(template.Must(template.ParseFiles(index,scripts,style,header,content,footer)))
		ui.SetInputChannel(make(chan []byte))
		ui.StartReceiver()
		
		env.Output.WriteChDebug("(UI::server::NewUI) New UI listening at: "+ui.Listenaddr)
	
		ui.clients = make( map[chan []byte]bool)
		ui.newClients = make( chan chan []byte)
		ui.defunctClients = make( chan chan []byte)

	}

	return ui
}

//#
//# Getters/Setters methods for Checks object
//#---------------------------------------------------------------------

//
//# SetListenaddr
func (u *UI) SetListenaddr(l string){
	env.Output.WriteChDebug("(UI::server::SetListenaddr) Set value")
	u.Listenaddr = l
}
//
//# SetRouter
func (u *UI) SetRouter(r *mux.Router){
	env.Output.WriteChDebug("(UI::server::SetRouter) Set value")
	u.router = r
}
//
//# SetTemplates
func (u *UI) SetTemplates(t *template.Template){
	env.Output.WriteChDebug("(UI::server::SetTemplates) Set value")
	u.templates = t
}
//
//# SetInputChannel
func (u *UI) SetInputChannel(i chan []byte){
	env.Output.WriteChDebug("(UI::server::SetInputChannel) Set value")
	u.inputChannel = i
}
//
//# GetListenaddr
func (u *UI) GetListenaddr() string {
	env.Output.WriteChDebug("(UI::server::GetListenaddr) Get value")
	return u.Listenaddr
}
//
//# GetRouter
func (u *UI) GetRouter() *mux.Router {
	env.Output.WriteChDebug("(UI::server::GetRouter) Get value")
	return u.router
}
//
//# GetTemplates
func (u *UI) GetTemplates() *template.Template {
	env.Output.WriteChDebug("(UI::server::GetTemplates) Get value")
	return u.templates
}
//
//# GetInputChannel
func (u *UI) GetInputChannel() chan []byte {
	env.Output.WriteChDebug("(UI::server::GetInputChannel) Get value")
	return u.inputChannel
}

//#
//# Specific methods
//#---------------------------------------------------------------------
//
//# SayHi: do nothing
func (u *UI) SayHi() {
  env.Output.WriteChInfo("(UI::server::SayHi) Hi! I'm your UI server instance.")
}
//
//# GetUI: method returns global ui
func GetUI() *UI {
	env.Output.WriteChDebug("(UI::server::GetUI) Get UI listening at: "+ui.Listenaddr)
	return ui
}
//
//# StartUI: method starts web server
func (u *UI) StartUI(){
	env.Output.WriteChDebug("(UI::server::StartUI) Starting UI listening at: "+u.Listenaddr)
	u.GenerateRoutes()
	u.router.Handle("/images/{img}",http.StripPrefix("/images/", http.FileServer(http.Dir("./ui/images/"))))
	u.router.Handle("/scripts/{script}",http.StripPrefix("/scripts/", http.FileServer(http.Dir("./ui/scripts/"))))
	u.router.Handle("/style/{style}",http.StripPrefix("/style/", http.FileServer(http.Dir("./ui/style/"))))
	
	log.Fatal(http.ListenAndServe(u.Listenaddr, u.router))
}

//
//# StartReceiver: method prepare engine to receive []byte to be sent to client
func (u *UI) StartReceiver() error {
	stormController := make(chan bool)
	enableDataReceiver := true

	if u == nil {
		return errors.New("(UI::server::StartReceiver) UI has not been initialized")
	}

	if u.inputChannel == nil {
		env.Output.WriteChDebug("(UI::server::StartReceiver) Initializing inputChannel")
		u.inputChannel = make(chan []byte)
	}

	stormControllerHandler := func () {
    env.Output.WriteChDebug("(UI::server::StartReceiver::stormController)")
    timeout := time.After(time.Duration(30) * time.Second)
    for{
      select{
      case <-timeout:
				stormController <- true
      }
    }
  }

	env.Output.WriteChDebug("(UI::server::StartReceiver) Starting byte receiver")
  go func() {
    defer close (u.inputChannel)
    for{
    	select{
    	// new client is connected
			case c := <-u.newClients:
				env.Output.WriteChDebug("(UI::server::StartReceiver) Add new client")
				u.clients[c] = true
			// client disconnected
			case c := <-u.defunctClients:
				env.Output.WriteChDebug("(UI::server::StartReceiver) Disconnection for client")
				delete(u.clients, c)
				close(c)
			// send data to clients
	    case data := <-u.inputChannel:
	      env.Output.WriteChDebug("(UI::server::StartReceiver) Data received")
		    if enableDataReceiver {
					for c, _ := range u.clients {
						c <- data
					}
		    	enableDataReceiver = false
		    	go stormControllerHandler()
		    } else {
		    	env.Output.WriteChDebug("(UI::server::StartReceiver) Data received is not enabled")
		    }
			case <- stormController:
				enableDataReceiver = true
				env.Output.WriteChDebug("(UI::server::StartReceiver) Data received enabled")
	    }
    }
  }()
  return nil
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