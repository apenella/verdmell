package ui

import(
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui", http.StatusFound)
}

func WebUI(w http.ResponseWriter, r *http.Request){
	env.Output.WriteChDebug("(UI::handler::WebUI)")

	if ui := GetUI(); ui != nil {
		// load template index.html
		if err := ui.templates.ExecuteTemplate(w,"index.html",nil); err != nil {
	 		http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		env.Output.WriteChError("(UI::handler::WebUI) The ui is null")		
	}
}

//#######################################################################################################