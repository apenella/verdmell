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
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, u *UI) error {
	http.Redirect(w, r, "/ui", http.StatusFound)
	return nil
}

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