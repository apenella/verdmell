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
	"fmt"
	"errors"
	"net/http"
)
//
//# Index
func Index(w http.ResponseWriter, r *http.Request, u *UI) error {
	http.Redirect(w, r, "/ui", http.StatusFound)
	return nil
}
//
//# WebUI
func WebUI(w http.ResponseWriter, r *http.Request, u *UI) error {
	env.Output.WriteChDebug("(UI::handler::WebUI)")

	if ui := GetUI(); ui != nil {
		// load template index.html
		if err := renderTemplate(w,"index.html",u); err != nil {
			return err
		}
	} else {
		msg := "(UI::handler::WebUI) UI has not been started yet"
		env.Output.WriteChError(msg)
		return errors.New(msg)		
	}
	return nil
}

func SSE(w http.ResponseWriter, r *http.Request, u *UI) error {
	env.Output.WriteChDebug("(UI::handler::SSE)")

	// Make sure that the writer supports flushing.
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return errors.New("(UI::handler::SSE) Streaming unsupported")
	}

	// Set the headers related to event streaming.
	//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	//CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	dataChan := make(chan []byte)
	u.newClients <- dataChan
	
	// Listen to the closing of the http connection via the CloseNotifier
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		// Remove this client from the map of attached clients
		// when `EventHandler` exits.
		u.defunctClients <- dataChan
		env.Output.WriteChDebug("(UI::handler::SSE) HTTP connection just closed.")
	}()

	for {
		// Read from our dataChan.
		data, open := <-dataChan
		if !open {
			env.Output.WriteChDebug("(UI::handler::SSE) Channel has been closed")
			break
		}
		env.Output.WriteChDebug("(UI::handler::SSE) Send event")
		// Write to the ResponseWriter, `w`.
		//w.Write(data)
		fmt.Fprintf(w, "data:%s\n\n", data)

		// Flush the response. This is only possible if
		// the repsonse supports streaming.
		f.Flush()			
	}

	env.Output.WriteChDebug("(UI::handler::SSE) Finished HTTP request at ", r.URL.Path)
	return nil
}
//#
//# Specific methods
//#---------------------------------------------------------------------

func renderTemplate(w http.ResponseWriter, tmpl string, u *UI) error {
  if err := ui.templates.ExecuteTemplate(w,tmpl,u); err != nil {
		return err
	}
	return nil
}

//#######################################################################################################